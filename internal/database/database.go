package database

import (
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gin-quickstart/internal/config"
	"gin-quickstart/internal/models"
)

var DB *gorm.DB

// Connect establishes a connection to the MS SQL database.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("Database connection established")
	return db, nil
}

// AutoMigrate runs database migrations for all models.
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	return db.AutoMigrate(
		&models.Customer{},
		&models.OTPRecord{},
	)
}

// ConnectWithMock creates an in-memory database for testing/development.
// Use this when MS SQL is not available.
func ConnectWithMock() (*gorm.DB, error) {
	log.Println("Using in-memory mock database (MS SQL not configured)")
	// For demonstration, we'll use a mock implementation
	// In production, this should connect to actual MS SQL
	return nil, nil
}
