package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PlaceOrder places a new order for the authenticated user
func CreateOrder(c *gin.Context) {
	var orderInput struct {
		Products []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required"`
		} `json:"products" binding:"required"`
	}
	if err := c.ShouldBindJSON(&orderInput); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	order := models.Order{
		UserID: userID.(uint),
		Status: "Pending",
	}
	if err := models.DB.Create(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	for _, product := range orderInput.Products {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		}
		if err := models.DB.Create(&orderItem).Error; err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to add products to order")
			return
		}
	}

	utils.RespondSuccess(c, "Order created successfully", order)
}

func PlaceOrder(c *gin.Context) {
	var orderInput struct {
		Products []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,min=1"`
		} `json:"products" binding:"required"`
		Total float64 `json:"total" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&orderInput); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid order data")
		return
	}

	// Get the authenticated user's ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Create the order in the database
	order := models.Order{
		UserID: userID.(uint),
		Total:  orderInput.Total,
		Status: "Pending",
	}
	if err := models.DB.Create(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// Add products to the order
	for _, product := range orderInput.Products {
		orderProduct := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		}
		if err := models.DB.Create(&orderProduct).Error; err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to add products to order")
			return
		}
	}

	utils.RespondSuccess(c, "Order placed successfully", order)
}

func GetOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var orders []models.Order
	if err := models.DB.Preload("Items").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	utils.RespondSuccess(c, "Orders fetched successfully", orders)
}

func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := models.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Order not found")
		return
	}

	if order.Status != "Pending" {
		utils.RespondError(c, http.StatusBadRequest, "Only pending orders can be cancelled")
		return
	}

	order.Status = "Cancelled"
	if err := models.DB.Save(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	utils.RespondSuccess(c, "Order cancelled successfully", nil)
}

// GetUserOrders fetches all orders for the authenticated user
func GetUserOrders(c *gin.Context) {
	// Get the authenticated user's ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var orders []models.Order
	if err := models.DB.Where("user_id = ?", userID).Preload("Products").Find(&orders).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	utils.RespondSuccess(c, "Orders retrieved successfully", orders)
}

// UpdateOrderStatus updates the status of an order (admin only)
func UpdateOrderStatus(c *gin.Context) {
	// Get order ID from the request URL
	orderID := c.Param("id")

	var order models.Order
	if err := models.DB.Preload("Products").First(&order, orderID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Order not found")
		return
	}

	// Check if the user is an admin
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		utils.RespondError(c, http.StatusForbidden, "Only admins can update order status")
		return
	}

	// Parse the new status from the request body
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Validate status
	allowedStatuses := map[string]bool{"Pending": true, "Shipped": true, "Cancelled": true}
	if !allowedStatuses[input.Status] {
		utils.RespondError(c, http.StatusBadRequest, "Invalid order status")
		return
	}

	order.Status = input.Status
	if err := models.DB.Save(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	utils.RespondSuccess(c, "Order status updated successfully", order)
}
