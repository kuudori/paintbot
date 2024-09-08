package middleware

import (
	"PaintBackend/internal/models"
	"PaintBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
)

func RequestProcessingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logRequest(c)
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered", "error", r, "path", c.Request.URL.Path)
				utils.Response(c, http.StatusOK, models.BaseResponse{Code: -999}) // Unhandled error
				c.Abort()
			}
		}()

		c.Next()
	}
}

func logRequest(c *gin.Context) {
	var bodyBytes []byte

	if c.ContentType() == "multipart/form-data" {
		bodyBytes = []byte("File multipart/form-data")
	} else {
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
	}

	slog.Info("Incoming request",
		"method", c.Request.Method,
		"url", c.Request.URL.String(),
		"headers", c.Request.Header,
		"body", string(bodyBytes),
	)
}

func ptr(s string) *string {
	return &s
}
