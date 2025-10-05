package middleware

import (
	"net/http"

	"docker-proxy/pkg/rules"

	"github.com/gin-gonic/gin"
)

// ACLMiddleware creates a middleware that enforces access control rules
func ACLMiddleware(matcher *rules.Matcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		if !matcher.IsAllowed(method, path) {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Access to this API endpoint is not allowed",
				"path":    path,
				"method":  method,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
