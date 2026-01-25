package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gin-quickstart/internal/models"
	"gin-quickstart/internal/services"
	"gin-quickstart/internal/utils"
)

// MockCustomerRepository is a mock implementation of CustomerRepositoryInterface
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(customer *models.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetAll() ([]models.Customer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindByID(id uuid.UUID) (*models.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindByPhone(phone string) (*models.Customer, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindByEmail(email string) (*models.Customer, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(customer *models.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockOTPRepository is a mock implementation of OTPRepositoryInterface
type MockOTPRepository struct {
	mock.Mock
}

func (m *MockOTPRepository) Create(otp *models.OTPRecord) error {
	args := m.Called(otp)
	return args.Error(0)
}

func (m *MockOTPRepository) FindLatestByPhone(phone string) (*models.OTPRecord, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.OTPRecord), args.Error(1)
}

func (m *MockOTPRepository) MarkAsUsed(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOTPRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func TestOnboardCustomer_Success(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	req := models.OnboardRequest{
		Phone:    "+2348012345678",
		Email:    "test@example.com",
		Password: "password123",
		State:    "Lagos",
		LGA:      "Ikeja",
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(nil, nil)
	mockCustomerRepo.On("FindByEmail", req.Email).Return(nil, nil)
	mockCustomerRepo.On("Create", mock.AnythingOfType("*models.Customer")).Return(nil)
	mockOTPRepo.On("Create", mock.AnythingOfType("*models.OTPRecord")).Return(nil)

	// Act
	err := service.OnboardCustomer(req)

	// Assert
	assert.NoError(t, err)
	mockCustomerRepo.AssertExpectations(t)
	mockOTPRepo.AssertExpectations(t)
}

func TestOnboardCustomer_InvalidLGA(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	req := models.OnboardRequest{
		Phone:    "+2348012345678",
		Email:    "test@example.com",
		Password: "password123",
		State:    "Lagos",
		LGA:      "InvalidLGA", // LGA doesn't exist in Lagos
	}

	// Act
	err := service.OnboardCustomer(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidLGA, err)
}

func TestOnboardCustomer_PhoneAlreadyExists(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	existingCustomer := &models.Customer{
		Phone: "+2348012345678",
	}

	req := models.OnboardRequest{
		Phone:    "+2348012345678",
		Email:    "test@example.com",
		Password: "password123",
		State:    "Lagos",
		LGA:      "Ikeja",
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(existingCustomer, nil)

	// Act
	err := service.OnboardCustomer(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrPhoneAlreadyExists, err)
}

func TestOnboardCustomer_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	existingCustomer := &models.Customer{
		Email: "test@example.com",
	}

	req := models.OnboardRequest{
		Phone:    "+2348012345678",
		Email:    "test@example.com",
		Password: "password123",
		State:    "Lagos",
		LGA:      "Ikeja",
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(nil, nil)
	mockCustomerRepo.On("FindByEmail", req.Email).Return(existingCustomer, nil)

	// Act
	err := service.OnboardCustomer(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrEmailAlreadyExists, err)
}

func TestVerifyOTP_Success(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	customer := &models.Customer{
		ID:              uuid.New(),
		Phone:           "+2348012345678",
		IsPhoneVerified: false,
	}

	otpRecord := utils.CreateOTPRecord(customer.Phone)

	req := models.VerifyOTPRequest{
		Phone: customer.Phone,
		OTP:   otpRecord.OTP,
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(customer, nil)
	mockOTPRepo.On("FindLatestByPhone", req.Phone).Return(otpRecord, nil)
	mockOTPRepo.On("MarkAsUsed", otpRecord.ID).Return(nil)
	mockCustomerRepo.On("Update", mock.AnythingOfType("*models.Customer")).Return(nil)

	// Act
	err := service.VerifyOTP(req)

	// Assert
	assert.NoError(t, err)
	mockCustomerRepo.AssertExpectations(t)
	mockOTPRepo.AssertExpectations(t)
}

func TestVerifyOTP_InvalidOTP(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	customer := &models.Customer{
		ID:              uuid.New(),
		Phone:           "+2348012345678",
		IsPhoneVerified: false,
	}

	otpRecord := utils.CreateOTPRecord(customer.Phone)

	req := models.VerifyOTPRequest{
		Phone: customer.Phone,
		OTP:   "000000", // Wrong OTP
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(customer, nil)
	mockOTPRepo.On("FindLatestByPhone", req.Phone).Return(otpRecord, nil)

	// Act
	err := service.VerifyOTP(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrInvalidOTP, err)
}

func TestVerifyOTP_PhoneNotRegistered(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	req := models.VerifyOTPRequest{
		Phone: "+2348099999999",
		OTP:   "123456",
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(nil, nil)

	// Act
	err := service.VerifyOTP(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrPhoneNotRegistered, err)
}

func TestVerifyOTP_AlreadyVerified(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	customer := &models.Customer{
		ID:              uuid.New(),
		Phone:           "+2348012345678",
		IsPhoneVerified: true, // Already verified
	}

	req := models.VerifyOTPRequest{
		Phone: customer.Phone,
		OTP:   "123456",
	}

	mockCustomerRepo.On("FindByPhone", req.Phone).Return(customer, nil)

	// Act
	err := service.VerifyOTP(req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, services.ErrPhoneAlreadyVerified, err)
}

func TestGetAll_Success(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	expectedCustomers := []models.Customer{
		{ID: uuid.New(), Phone: "+2348012345678", Email: "test1@example.com"},
		{ID: uuid.New(), Phone: "+2348087654321", Email: "test2@example.com"},
	}

	mockCustomerRepo.On("GetAll").Return(expectedCustomers, nil)

	// Act
	customers, err := service.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, customers, 2)
	mockCustomerRepo.AssertExpectations(t)
}

func TestResendOTP_Success(t *testing.T) {
	// Arrange
	mockCustomerRepo := new(MockCustomerRepository)
	mockOTPRepo := new(MockOTPRepository)
	service := services.NewCustomerService(mockCustomerRepo, mockOTPRepo)

	customer := &models.Customer{
		ID:              uuid.New(),
		Phone:           "+2348012345678",
		IsPhoneVerified: false,
	}

	mockCustomerRepo.On("FindByPhone", customer.Phone).Return(customer, nil)
	mockOTPRepo.On("Create", mock.AnythingOfType("*models.OTPRecord")).Return(nil)

	// Act
	err := service.ResendOTP(customer.Phone)

	// Assert
	assert.NoError(t, err)
	mockCustomerRepo.AssertExpectations(t)
	mockOTPRepo.AssertExpectations(t)
}
