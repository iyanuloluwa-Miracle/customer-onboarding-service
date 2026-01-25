package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/internal/utils"
)

// Common errors
var (
	ErrPhoneAlreadyExists  = errors.New("phone number already registered")
	ErrEmailAlreadyExists  = errors.New("email already registered")
	ErrInvalidLGA          = errors.New("LGA does not belong to selected state")
	ErrInvalidOTP          = errors.New("invalid or expired OTP")
	ErrCustomerNotFound    = errors.New("customer not found")
	ErrPasswordHashFailed  = errors.New("failed to hash password")
	ErrPhoneNotRegistered  = errors.New("phone number not registered")
	ErrPhoneAlreadyVerified = errors.New("phone number already verified")
)

// CustomerServiceInterface defines the contract for customer business logic.
type CustomerServiceInterface interface {
	OnboardCustomer(req models.OnboardRequest) error
	VerifyOTP(req models.VerifyOTPRequest) error
	ResendOTP(phone string) error
	GetAll() ([]models.Customer, error)
	GetByPhone(phone string) (*models.Customer, error)
}

// CustomerService handles customer business logic.
type CustomerService struct {
	customerRepo repositories.CustomerRepositoryInterface
	otpRepo      repositories.OTPRepositoryInterface
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(
	customerRepo repositories.CustomerRepositoryInterface,
	otpRepo repositories.OTPRepositoryInterface,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		otpRepo:      otpRepo,
	}
}

// OnboardCustomer handles the customer onboarding process.
// Business Rules:
// 1. Phone number must be unique
// 2. Email must be unique
// 3. LGA must belong to the selected state
// 4. OTP is sent to verify phone number
func (s *CustomerService) OnboardCustomer(req models.OnboardRequest) error {
	// Validate LGA belongs to state
	if !utils.IsValidLGA(req.State, req.LGA) {
		return ErrInvalidLGA
	}

	// Check if phone already exists
	existingByPhone, _ := s.customerRepo.FindByPhone(req.Phone)
	if existingByPhone != nil {
		return ErrPhoneAlreadyExists
	}

	// Check if email already exists
	existingByEmail, _ := s.customerRepo.FindByEmail(req.Email)
	if existingByEmail != nil {
		return ErrEmailAlreadyExists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHashFailed
	}

	// Create customer (phone not yet verified)
	customer := &models.Customer{
		Phone:           req.Phone,
		Email:           req.Email,
		PasswordHash:    string(hash),
		State:           req.State,
		LGA:             req.LGA,
		IsPhoneVerified: false,
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return err
	}

	// Generate and send OTP
	otpRecord := utils.CreateOTPRecord(req.Phone)
	if err := s.otpRepo.Create(otpRecord); err != nil {
		return err
	}

	// Mock send OTP
	utils.MockSendOTP(req.Phone, otpRecord.OTP)

	return nil
}

// VerifyOTP verifies the OTP and completes onboarding.
func (s *CustomerService) VerifyOTP(req models.VerifyOTPRequest) error {
	// Find customer by phone
	customer, err := s.customerRepo.FindByPhone(req.Phone)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrPhoneNotRegistered
	}

	// Check if already verified
	if customer.IsPhoneVerified {
		return ErrPhoneAlreadyVerified
	}

	// Get latest OTP for this phone
	otpRecord, err := s.otpRepo.FindLatestByPhone(req.Phone)
	if err != nil {
		return err
	}

	// Validate OTP
	if !utils.ValidateOTP(otpRecord, req.OTP) {
		return ErrInvalidOTP
	}

	// Mark OTP as used
	if err := s.otpRepo.MarkAsUsed(otpRecord.ID); err != nil {
		return err
	}

	// Mark phone as verified
	customer.IsPhoneVerified = true
	return s.customerRepo.Update(customer)
}

// ResendOTP generates and sends a new OTP.
func (s *CustomerService) ResendOTP(phone string) error {
	// Find customer by phone
	customer, err := s.customerRepo.FindByPhone(phone)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrPhoneNotRegistered
	}

	// Check if already verified
	if customer.IsPhoneVerified {
		return ErrPhoneAlreadyVerified
	}

	// Generate and send new OTP
	otpRecord := utils.CreateOTPRecord(phone)
	if err := s.otpRepo.Create(otpRecord); err != nil {
		return err
	}

	utils.MockSendOTP(phone, otpRecord.OTP)
	return nil
}

// GetAll retrieves all customers.
func (s *CustomerService) GetAll() ([]models.Customer, error) {
	return s.customerRepo.GetAll()
}

// GetByPhone retrieves a customer by phone number.
func (s *CustomerService) GetByPhone(phone string) (*models.Customer, error) {
	return s.customerRepo.FindByPhone(phone)
}
