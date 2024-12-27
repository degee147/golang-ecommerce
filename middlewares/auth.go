package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", authHeader) // Debug log

		// Check if Authorization header is missing or does not start with "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract JWT secret key from environment variables
		key := os.Getenv("JWT_SECRET")
		if key == "" {
			fmt.Println("Error: JWT_SECRET environment variable is missing.") // Debug log
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix from the token
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Token String:", tokenStr) // Debug log

		// Parse the token using the secret key
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Check if the token's signing method is valid
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(key), nil
		})

		// Handle token parsing errors
		if err != nil {
			fmt.Println("Token Parse Error:", err) // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user info from the token (assuming `is_admin` is part of the token claims)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set the user ID and is_admin in the context
			c.Set("user_id", claims["user_id"])
			c.Set("is_admin", claims["is_admin"]) // Set is_admin value in context

			// Proceed to the next handler
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
		// Proceed to the next middleware or handler
		c.Next()
	}
}
