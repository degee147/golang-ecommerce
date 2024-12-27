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

	// Check if a product with the same name already exists
	var existingProduct models.Product
	if err := models.DB.Where("name = ?", product.Name).First(&existingProduct).Error; err == nil {
		// Product with the same name already exists
		utils.RespondError(c, http.StatusConflict, "Product with this name already exists")
		return
	}

	// Create the new product if no conflict
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
	var input struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Description string  `json:"description"`
	}

	// Bind input JSON to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	// Get product ID from URL parameter
	productID := c.Param("id")

	// Find the product by ID
	var product models.Product
	if err := models.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Product not found")
		return
	}

	// Update product fields
	product.Name = input.Name
	product.Price = input.Price
	product.Stock = input.Stock
	product.Description = input.Description

	// Save the updated product to the database
	if err := models.DB.Save(&product).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update product")
		return
	}

	utils.RespondSuccess(c, "Product updated successfully", product)
}

// DeleteProduct deletes a product (admin only)
func DeleteProduct(c *gin.Context) {
	// Get product ID from URL parameter
	productID := c.Param("id")

	// Find the product by ID
	var product models.Product
	if err := models.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Product not found")
		return
	}

	// Delete the product from the database
	if err := models.DB.Delete(&product).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	utils.RespondSuccess(c, "Product deleted successfully", nil)
}
