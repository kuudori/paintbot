package validators

import (
	"github.com/gin-gonic/gin"
)

func ValidateRequestBody(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBind(obj); err != nil {
		return err
	}

	return nil
}
