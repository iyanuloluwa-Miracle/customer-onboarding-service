package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gin-quickstart/internal/models"
	"gin-quickstart/internal/utils"
)

func TestGenerateOTP(t *testing.T) {
	// Act
	otp := utils.GenerateOTP()

	// Assert
	assert.Len(t, otp, 6, "OTP should be 6 digits")
	assert.Regexp(t, `^\d{6}$`, otp, "OTP should contain only digits")
}

func TestGenerateOTP_Uniqueness(t *testing.T) {
	// Generate multiple OTPs and check they're not all the same
	otps := make(map[string]bool)
	for i := 0; i < 100; i++ {
		otp := utils.GenerateOTP()
		otps[otp] = true
	}

	// With 100 random 6-digit numbers, we should have at least some unique values
	assert.Greater(t, len(otps), 1, "OTPs should be unique")
}

func TestCreateOTPRecord(t *testing.T) {
	// Arrange
	phone := "+2348012345678"

	// Act
	record := utils.CreateOTPRecord(phone)

	// Assert
	assert.Equal(t, phone, record.Phone)
	assert.Len(t, record.OTP, 6)
	assert.False(t, record.IsUsed)
	assert.True(t, record.ExpiresAt.After(time.Now()))
}

func TestValidateOTP_Success(t *testing.T) {
	// Arrange
	record := utils.CreateOTPRecord("+2348012345678")

	// Act
	result := utils.ValidateOTP(record, record.OTP)

	// Assert
	assert.True(t, result)
}

func TestValidateOTP_NilRecord(t *testing.T) {
	// Act
	result := utils.ValidateOTP(nil, "123456")

	// Assert
	assert.False(t, result)
}

func TestValidateOTP_WrongOTP(t *testing.T) {
	// Arrange
	record := utils.CreateOTPRecord("+2348012345678")

	// Act
	result := utils.ValidateOTP(record, "000000")

	// Assert
	assert.False(t, result)
}

func TestValidateOTP_UsedOTP(t *testing.T) {
	// Arrange
	record := utils.CreateOTPRecord("+2348012345678")
	record.IsUsed = true

	// Act
	result := utils.ValidateOTP(record, record.OTP)

	// Assert
	assert.False(t, result)
}

func TestValidateOTP_ExpiredOTP(t *testing.T) {
	// Arrange
	record := &models.OTPRecord{
		Phone:     "+2348012345678",
		OTP:       "123456",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		IsUsed:    false,
	}

	// Act
	result := utils.ValidateOTP(record, record.OTP)

	// Assert
	assert.False(t, result)
}

func TestIsValidLGA_ValidCombination(t *testing.T) {
	testCases := []struct {
		state string
		lga   string
	}{
		{"Lagos", "Ikeja"},
		{"Lagos", "Epe"},
		{"Lagos", "Ikorodu"},
		{"Oyo", "Ibadan North"},
		{"Oyo", "Ibadan South"},
	}

	for _, tc := range testCases {
		t.Run(tc.state+"_"+tc.lga, func(t *testing.T) {
			result := utils.IsValidLGA(tc.state, tc.lga)
			assert.True(t, result)
		})
	}
}

func TestIsValidLGA_InvalidCombination(t *testing.T) {
	testCases := []struct {
		state string
		lga   string
	}{
		{"Lagos", "Ibadan North"}, // Ibadan North is in Oyo, not Lagos
		{"Oyo", "Ikeja"},          // Ikeja is in Lagos, not Oyo
		{"Kano", "Ikeja"},         // Kano is not in the map
		{"Lagos", "InvalidLGA"},   // LGA doesn't exist
	}

	for _, tc := range testCases {
		t.Run(tc.state+"_"+tc.lga, func(t *testing.T) {
			result := utils.IsValidLGA(tc.state, tc.lga)
			assert.False(t, result)
		})
	}
}
