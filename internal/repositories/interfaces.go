package repositories

import (
	"gin-quickstart/internal/models"

	"github.com/google/uuid"
)

// CustomerRepositoryInterface defines the contract for customer data operations.
// This interface allows for easy mocking in unit tests.
type CustomerRepositoryInterface interface {
	Create(customer *models.Customer) error
	GetAll() ([]models.Customer, error)
	FindByID(id uuid.UUID) (*models.Customer, error)
	FindByPhone(phone string) (*models.Customer, error)
	FindByEmail(email string) (*models.Customer, error)
	Update(customer *models.Customer) error
	Delete(id uuid.UUID) error
}

// OTPRepositoryInterface defines the contract for OTP data operations.
type OTPRepositoryInterface interface {
	Create(otp *models.OTPRecord) error
	FindLatestByPhone(phone string) (*models.OTPRecord, error)
	MarkAsUsed(id uuid.UUID) error
	DeleteExpired() error
}
