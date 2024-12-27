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

		// Validate if the token is valid and extract the claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Example of extracting a user ID from the claims
			userID, ok := claims["user_id"].(float64) // Assuming user_id is a number
			if !ok {
				fmt.Println("Error: Invalid user ID in token") // Debug log
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			// Set the user ID in the context for further use in the application
			c.Set("user_id", uint(userID))
			fmt.Println("User ID from token:", userID) // Debug log
		} else {
			fmt.Println("Error: Invalid token claims") // Debug log
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Proceed to the next middleware or handler
		c.Next()
	}
}
