package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateProduct creates a new product (admin only)
// @Summary Create a new product
// @Description Creates a new product, requires admin access
// @Tags Products
// @Accept  json
// @Produce  json
// @Param product body models.Product true "Product data"
// @Success 201 {object} models.Product "Product created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 409 {object} utils.ErrorResponse "Product with this name already exists"
// @Failure 500 {object} utils.ErrorResponse "Failed to create product"
// @Security BearerAuth
// @Router /products [post]
//
//	@Example {
//	   "name": "Product 3",
//	   "price": 75,
//	   "stock": 30,
//	   "description": "A description of the product 3"
//	}
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
// @Summary Get all products
// @Description Fetches all products in the system
// @Tags Products
// @Produce  json
// @Success 200 {array} models.Product "List of products"
// @Failure 500 {object} utils.ErrorResponse "Failed to fetch products"
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := models.DB.Find(&products).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	utils.RespondSuccess(c, "Products fetched successfully", products)
}

// UpdateProduct updates an existing product (admin only)
// @Summary Update an existing product
// @Description Updates the details of an existing product, requires admin access
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product data"
// @Success 200 {object} models.Product "Product updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid input data"
// @Failure 404 {object} utils.ErrorResponse "Product not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to update product"
// @Security BearerAuth
// @Router /products/{id} [put]
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
// @Summary Delete a product
// @Description Deletes a product by its ID, requires admin access
// @Tags Products
// @Param id path string true "Product ID"
// @Success 200 {string} string "Product deleted successfully"
// @Failure 404 {object} utils.ErrorResponse "Product not found"
// @Failure 500 {object} utils.ErrorResponse "Failed to delete product"
// @Security BearerAuth
// @Router /products/{id} [delete]
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
