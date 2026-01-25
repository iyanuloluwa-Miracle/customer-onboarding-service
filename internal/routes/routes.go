package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"gin-quickstart/internal/clients"
	"gin-quickstart/internal/controller"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/internal/services"
)

// SetupRoutes configures all API routes.
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	customerRepo := repositories.NewCustomerRepository(db)
	otpRepo := repositories.NewOTPRepository(db)

	// Initialize services
	customerService := services.NewCustomerService(customerRepo, otpRepo)

	// Initialize clients
	bankClient := clients.NewBankClient()

	// Initialize controllers
	customerController := controller.NewCustomerController(customerService, bankClient)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "pong",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Customer routes
		customers := v1.Group("/customers")
		{
			customers.POST("/onboard", customerController.Onboard)
			customers.POST("/verify-otp", customerController.VerifyOTP)
			customers.POST("/resend-otp", customerController.ResendOTP)
			customers.GET("", customerController.GetAll)
		}

		// Bank routes
		v1.GET("/banks", customerController.GetBanks)
	}
}
