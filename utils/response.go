package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondSuccess sends a success response
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// RespondError sends an error response
func RespondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}

// RespondValidationError sends a validation error response
func RespondValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": "Validation failed",
		"errors":  errors,
	})
}

type ErrorResponse struct {
	Message string `json:"message"` // Message field to hold the error message
}
type SuccessResponse struct {
	Message string `json:"message"` // Message field to hold the error message
}

type LoginResponse struct {
	Token string `json:"token"`
}
