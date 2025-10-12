package middleware

import (
	"net/http"

	"dockershield/pkg/rules"

	"github.com/gin-gonic/gin"
)

// ACLMiddleware creates a middleware that enforces access control rules
func ACLMiddleware(matcher *rules.Matcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if advanced filter already authorized this request
		// This allows DKRPRX__ variables to override ACL settings
		if authorized, exists := c.Get("advanced_filter_authorized"); exists && authorized.(bool) {
			c.Next()
			return
		}

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
