package config

import (
	"os"
)

// Config holds application configuration.
type Config struct {
	// Server configuration
	ServerPort string
	GinMode    string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// OTP configuration
	OTPExpiryMinutes int
}

// Load returns the application configuration.
// Values are loaded from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		// Server
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "1433"),
		DBUser:     getEnv("DB_USER", "sa"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "customer_onboarding"),

		// OTP
		OTPExpiryMinutes: 5,
	}
}

// GetDSN returns the database connection string for MS SQL Server.
func (c *Config) GetDSN() string {
	return "sqlserver://" + c.DBUser + ":" + c.DBPassword +
		"@" + c.DBHost + ":" + c.DBPort +
		"?database=" + c.DBName
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
