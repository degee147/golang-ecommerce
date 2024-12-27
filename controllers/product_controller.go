package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateProduct creates a new product (admin only)
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := models.DB.Create(&product).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create product")
		return
	}

	utils.RespondSuccess(c, "Product created successfully", product)
}

// GetProducts fetches all products
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := models.DB.Find(&products).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	utils.RespondSuccess(c, "Products fetched successfully", products)
}

// UpdateProduct updates a product (admin only)
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := models.DB.Where("id = ?", id).First(&product).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Product not found")
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := models.DB.Save(&product).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update product")
		return
	}

	utils.RespondSuccess(c, "Product updated successfully", product)
}

// DeleteProduct deletes a product (admin only)
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Delete(&models.Product{}, id).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	utils.RespondSuccess(c, "Product deleted successfully", nil)
}
