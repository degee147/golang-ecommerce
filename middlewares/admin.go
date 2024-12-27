package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the `is_admin` value set in the JWT middleware or context
		isAdmin, exists := c.Get("is_admin")

		// Check if the value exists and is a boolean
		if !exists || isAdmin == nil || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden: Admins only"})
			c.Abort() // Stop the request from proceeding further
			return
		}

		// If the user is an admin, continue with the request
		c.Next()
	}
}
