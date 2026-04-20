package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/SuperJinggg/ai-router/internal/common"
	"github.com/SuperJinggg/ai-router/internal/errno"
	"github.com/SuperJinggg/ai-router/internal/service"
)

func IPBlacklistFilter(blacklistService *service.BlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if blacklistService.IsBlocked(clientIP) {
			log.Printf("blocked request from blacklisted IP: ip=%s path=%s method=%s", clientIP, c.Request.URL.Path, c.Request.Method)
			c.JSON(http.StatusForbidden, common.BaseResponse{
				Code:    errno.ForbiddenError.Code,
				Data:    nil,
				Message: "您的 IP 已被封禁",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
