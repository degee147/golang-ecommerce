package routes

import (
	"ecommerce-api/controllers"
	"ecommerce-api/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// User routes
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.Login)

		// Product routes
		api.GET("/products", controllers.GetProducts)
		api.POST("/products", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.CreateProduct)
		api.PUT("/products/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.UpdateProduct)
		api.DELETE("/products/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.DeleteProduct)

		// Order routes
		api.POST("/orders", middlewares.AuthMiddleware(), controllers.CreateOrder)
		api.GET("/orders", middlewares.AuthMiddleware(), controllers.GetOrders)
		api.PUT("/orders/:id", middlewares.AuthMiddleware(), middlewares.AdminMiddleware(), controllers.UpdateOrderStatus)
		api.DELETE("/orders/:id", middlewares.AuthMiddleware(), controllers.CancelOrder)
	}
}
