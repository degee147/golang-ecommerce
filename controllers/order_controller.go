package controllers

import (
	"ecommerce-api/models"
	"ecommerce-api/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PlaceOrder places a new order for the authenticated user
func CreateOrder(c *gin.Context) {
	// Get user ID from the context (set by the AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized user not found")
		return
	}

	// Convert userID from float64 to uint
	userIDUint, ok := userID.(float64)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	convertedUserID := uint(userIDUint)

	// Define the structure for the incoming order request
	var orderInput struct {
		Items []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  uint `json:"quantity" binding:"required"`
		} `json:"items" binding:"required"`
		Total float64 `json:"total" binding:"required"`
	}

	// Bind the incoming JSON to the orderInput struct
	if err := c.ShouldBindJSON(&orderInput); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid order data")
		return
	}

	// Create a new order record
	order := models.Order{
		UserID: convertedUserID,
		Total:  orderInput.Total,
		Status: "Pending",
	}

	// Save the order to the database
	if err := models.DB.Create(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// Check if products exist before adding them to the order
	for _, item := range orderInput.Items {
		var product models.Product
		// Check if the product exists
		if err := models.DB.First(&product, item.ProductID).Error; err != nil {
			utils.RespondError(c, http.StatusNotFound, fmt.Sprintf("Product with ID %d not found", item.ProductID))
			return
		}

		// Convert uint to int for Quantity field
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity), // Convert to int here
		}

		// Save the order items to the database
		if err := models.DB.Create(&orderItem).Error; err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Failed to add items to order")
			return
		}
	}

	// Respond with the created order details
	utils.RespondSuccess(c, "Order created successfully", order)
}

func GetOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	fmt.Println("userID String:", userID) // Debug log
	fmt.Println("exists :", exists)       // Debug log
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized user not found")
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
	// Retrieve the order ID from the URL
	orderID := c.Param("id")
	if orderID == "" {
		utils.RespondError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	// Retrieve the user ID from the context (set by the AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "Unauthorized user not found")
		return
	}

	// Convert userID from float64 to uint
	userIDUint, ok := userID.(float64)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	convertedUserID := uint(userIDUint)

	// Find the order by ID
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, fmt.Sprintf("Order with ID %s not found", orderID))
		return
	}

	// Check if the order belongs to the authenticated user
	if order.UserID != convertedUserID {
		utils.RespondError(c, http.StatusForbidden, "You do not have permission to cancel this order")
		return
	}

	// Check if the order status is already "Canceled" or cannot be canceled
	if order.Status == "Canceled" {
		utils.RespondError(c, http.StatusBadRequest, "Order is already canceled")
		return
	}

	// Update the order status to "Canceled"
	order.Status = "Canceled"
	if err := models.DB.Save(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to cancel the order")
		return
	}

	// Respond with the updated order details
	utils.RespondSuccess(c, "Order canceled successfully", order)
}

// GetUserOrders fetches all orders for the authenticated user
func GetUserOrders(c *gin.Context) {
	// Get the authenticated user's ID from the context
	userID, exists := c.Get("user_id")
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
	// Retrieve the order ID from the URL
	orderID := c.Param("id")
	if orderID == "" {
		utils.RespondError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	// Define the structure for the incoming request data
	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	// Bind the incoming JSON to the statusUpdate struct
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Find the order by ID
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, fmt.Sprintf("Order with ID %s not found", orderID))
		return
	}

	// Ensure that only valid statuses are accepted (e.g., Pending, Shipped, Delivered, Canceled)
	validStatuses := map[string]bool{
		"Pending":   true,
		"Shipped":   true,
		"Delivered": true,
		"Canceled":  true,
	}

	if _, valid := validStatuses[statusUpdate.Status]; !valid {
		utils.RespondError(c, http.StatusBadRequest, "Invalid status value")
		return
	}

	// Update the order status
	order.Status = statusUpdate.Status
	if err := models.DB.Save(&order).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to update order status")
		return
	}

	// Respond with the updated order details
	utils.RespondSuccess(c, "Order status updated successfully", order)
}
