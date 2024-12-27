package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	if err := models.DB.Create(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	utils.RespondSuccess(c, "User registered successfully", nil)
}

// Login and generate JWT
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Verify password
	if err := utils.VerifyPassword(user.Password, input.Password); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.RespondSuccess(c, "Login successful", gin.H{"token": token})
}
