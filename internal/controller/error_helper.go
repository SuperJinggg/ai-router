package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SuperJinggg/ai-router/internal/common"
	"github.com/SuperJinggg/ai-router/internal/errno"
)

func writeError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, common.BaseResponse{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}

func writeServiceError(ctx *gin.Context, err error) {
	if bizErr, ok := errno.AsBusinessError(err); ok {
		writeError(ctx, bizErr.Code, bizErr.Message)
		return
	}
	writeError(ctx, errno.SystemError.Code, "系统错误")
}
