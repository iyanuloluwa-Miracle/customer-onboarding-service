package repositories_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
)

func TestCustomerRepository_Create(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil) // In-memory mode
	customer := &models.Customer{
		Phone: "+2348012345678",
		Email: "test@example.com",
		State: "Lagos",
		LGA:   "Ikeja",
	}

	// Act
	err := repo.Create(customer)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, customer.ID)
}

func TestCustomerRepository_GetAll(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil)
	customer1 := &models.Customer{Phone: "+2348012345678", Email: "test1@example.com"}
	customer2 := &models.Customer{Phone: "+2348087654321", Email: "test2@example.com"}
	_ = repo.Create(customer1)
	_ = repo.Create(customer2)

	// Act
	customers, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, customers, 2)
}

func TestCustomerRepository_FindByPhone(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil)
	customer := &models.Customer{Phone: "+2348012345678", Email: "test@example.com"}
	_ = repo.Create(customer)

	// Act
	found, err := repo.FindByPhone("+2348012345678")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, customer.Phone, found.Phone)
}

func TestCustomerRepository_FindByPhone_NotFound(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil)

	// Act
	found, err := repo.FindByPhone("+2348099999999")

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestCustomerRepository_FindByEmail(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil)
	customer := &models.Customer{Phone: "+2348012345678", Email: "test@example.com"}
	_ = repo.Create(customer)

	// Act
	found, err := repo.FindByEmail("test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, customer.Email, found.Email)
}

func TestCustomerRepository_Update(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil)
	customer := &models.Customer{Phone: "+2348012345678", Email: "test@example.com"}
	_ = repo.Create(customer)

	// Act
	customer.IsPhoneVerified = true
	err := repo.Update(customer)

	// Assert
	assert.NoError(t, err)

	// Verify update
	found, _ := repo.FindByPhone(customer.Phone)
	assert.True(t, found.IsPhoneVerified)
}

func TestCustomerRepository_Delete(t *testing.T) {
	// Arrange
	repo := repositories.NewCustomerRepository(nil)
	customer := &models.Customer{Phone: "+2348012345678", Email: "test@example.com"}
	_ = repo.Create(customer)

	// Act
	err := repo.Delete(customer.ID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	customers, _ := repo.GetAll()
	assert.Len(t, customers, 0)
}

func TestOTPRepository_Create(t *testing.T) {
	// Arrange
	repo := repositories.NewOTPRepository(nil)
	otp := &models.OTPRecord{
		Phone:     "+2348012345678",
		OTP:       "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Act
	err := repo.Create(otp)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, otp.ID)
}

func TestOTPRepository_FindLatestByPhone(t *testing.T) {
	// Arrange
	repo := repositories.NewOTPRepository(nil)
	phone := "+2348012345678"
	
	// Create OTPs with explicit timestamps
	now := time.Now()
	otp1 := &models.OTPRecord{
		Phone:     phone,
		OTP:       "111111",
		ExpiresAt: now.Add(5 * time.Minute),
		CreatedAt: now.Add(-2 * time.Minute), // Older
	}
	otp2 := &models.OTPRecord{
		Phone:     phone,
		OTP:       "222222",
		ExpiresAt: now.Add(5 * time.Minute),
		CreatedAt: now.Add(-1 * time.Minute), // Newer
	}
	_ = repo.Create(otp1)
	_ = repo.Create(otp2)

	// Act
	found, err := repo.FindLatestByPhone(phone)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "222222", found.OTP) // Should be the latest
}

func TestOTPRepository_MarkAsUsed(t *testing.T) {
	// Arrange
	repo := repositories.NewOTPRepository(nil)
	otp := &models.OTPRecord{
		Phone:     "+2348012345678",
		OTP:       "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	_ = repo.Create(otp)

	// Act
	err := repo.MarkAsUsed(otp.ID)

	// Assert
	assert.NoError(t, err)

	// Should not find it anymore (it's used)
	found, _ := repo.FindLatestByPhone(otp.Phone)
	assert.Nil(t, found)
}

func TestOTPRepository_FindLatestByPhone_ExpiredExcluded(t *testing.T) {
	// Arrange
	repo := repositories.NewOTPRepository(nil)
	phone := "+2348012345678"
	expiredOTP := &models.OTPRecord{
		Phone:     phone,
		OTP:       "111111",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	}
	_ = repo.Create(expiredOTP)

	// Act
	found, err := repo.FindLatestByPhone(phone)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, found) // Expired OTP should not be found
}
