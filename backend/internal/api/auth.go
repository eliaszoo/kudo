package api

import (
	"net/http"
	"reward-system/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for wechat webhook (handled separately)
		if c.Request.URL.Path == "/api/v1/wechat" {
			c.Next()
			return
		}

		// Bearer token auth for API
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Missing authorization header"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid authorization header format"})
			c.Abort()
			return
		}

		if tokenParts[1] != cfg.APIToken {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}