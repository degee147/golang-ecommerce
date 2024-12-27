package main

import (
	"ecommerce-api/models"
	"ecommerce-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	models.ConnectDatabase()

	// Auto-migrate database models
	models.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})

	// Set up routes
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Use(cors.Default())

	// Start server
	router.Run(":8080") // Run on port 8080
}
