package api

import (
	"PaintBackend/internal/models"
	database "PaintBackend/internal/storage/models"
	"PaintBackend/internal/utils"
	"PaintBackend/internal/validators"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	filepathpackage "path/filepath"
	"strconv"
)

type Handlers struct {
	FileStorage database.FileRepository
}

func NewHandler(fileStorage database.FileRepository) *Handlers {
	return &Handlers{FileStorage: fileStorage}
}

func (h *Handlers) UploadImage(ctx *gin.Context) {
	user, err := validators.ValidateWebAppData(ctx.GetHeader("X-Auth"))
	if err != nil {
		utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: -1}) // Unauthorized
		ctx.Abort()
		return
	}
	var body models.FileUpload
	var message string
	if err := validators.ValidateRequestBody(ctx, &body); err != nil {
		exc := err.Error()
		utils.Response(ctx, http.StatusBadRequest, models.BaseResponse{Code: 1001, MessageDetails: &exc})
		return
	}
	file, _ := ctx.FormFile("image")

	if !validators.IsJPG(file) {
		messageDetails := "Wrong file format"
		utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1001, MessageDetails: &messageDetails})
		return
	}

	// save file
	uploadDir := filepathpackage.Join("media", fmt.Sprintf("%v", user.Id))
	imageURL, filepath, err := utils.SaveFile(uploadDir, file)
	if err != nil {
		slog.Error("SaveFile error: ", err)
		utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1003}) // Failed to upload image
		return
	}
	newFile := database.DbFile{
		ChatID:   user.Id,
		FileUrl:  imageURL,
		FilePath: filepath,
	}
	err = h.FileStorage.Create(ctx, &newFile)
	if err != nil {
		slog.Error("AddFileToDB error: ", err)

		utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1003}) // Upload error
		return
	}
	message = "Image uploaded successfully"
	utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 0, Message: &message})

}

// DeleteImage godoc
// @Summary      Delete an image
// @Description  Delete an image file by id
// @Tags         images
// @Accept       json
// @Produce      json
// @Param        id path string true "id of the image to delete"
// @Router       /images/{id} [delete]
func (h *Handlers) DeleteImage(ctx *gin.Context) {

	fileId, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)
	file, err := h.FileStorage.GetFileById(ctx, fileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1005}) // File not found
		} else {
			utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1006}) // Database error
		}
		return
	}

	_, err = h.FileStorage.GetFileById(ctx, fileId)
	if err != nil {
		utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1006}) // Database error
		return
	}

	err = os.Remove(file.FilePath)
	if err != nil {
		utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 1007}) // Database error
		return
	}

	message := "Image deleted successfully"
	utils.Response(ctx, http.StatusOK, models.BaseResponse{Code: 0, Message: &message})
}
