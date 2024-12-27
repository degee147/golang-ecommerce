package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser registers a new user in the system
// @Summary Register a new user
// @Description Register a new user with email, password, and admin status
// @Tags Users
// @Accept  json
// @Produce  json
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 500 {object} utils.ErrorResponse "Failed to create user"
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		IsAdmin  bool   `json:"is_admin" binding:"required"`
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
		IsAdmin:  input.IsAdmin,
	}
	if err := models.DB.Create(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	utils.RespondSuccess(c, "User registered successfully", nil)
}

// Login authenticates a user and generates a JWT
// @Summary Login and get JWT token
// @Description Authenticate a user with email and password, then generate a JWT token
// @Tags Users
// @Accept  json
// @Produce  json
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 401 {object} utils.ErrorResponse "Invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Failed to generate token"
// @Router /login [post]
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
