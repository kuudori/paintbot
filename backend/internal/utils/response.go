package utils

import (
	"PaintBackend/constants"
	"PaintBackend/internal/models"
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, status int, responseModel models.BaseResponse) {
	if responseModel.Code != 0 {
		if msg, ok := constants.ErrorMessages[responseModel.Code]; ok {
			responseModel.Message = &msg
			c.AbortWithStatusJSON(status, responseModel)
			return
		}
	}
	c.JSON(status, responseModel)
}
