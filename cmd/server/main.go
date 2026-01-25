package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gin-quickstart/internal/config"
	"gin-quickstart/internal/database"
	"gin-quickstart/internal/routes"

	_ "gin-quickstart/docs" // Swagger docs
)

// @title Customer Onboarding API
// @version 1.0
// @description A microservice API for customer onboarding with phone verification via OTP.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @schemes http https

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Connect to database
	var db *gorm.DB
	var err error

	if cfg.DBPassword != "" {
		db, err = database.Connect(cfg)
		if err != nil {
			log.Printf("Warning: Could not connect to MS SQL database: %v", err)
			log.Println("Falling back to in-memory storage")
			db = nil
		} else {
			// Run migrations
			if err := database.AutoMigrate(db); err != nil {
				log.Fatalf("Failed to run migrations: %v", err)
			}
		}
	} else {
		log.Println("DB_PASSWORD not set. Using in-memory storage.")
		log.Println("Set DB_HOST, DB_USER, DB_PASSWORD, DB_NAME environment variables for MS SQL.")
	}

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	log.Printf("Swagger UI available at http://localhost%s/swagger/index.html", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}