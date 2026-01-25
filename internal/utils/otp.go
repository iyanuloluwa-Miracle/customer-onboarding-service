package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"gin-quickstart/internal/models"
)

const (
	// OTPLength defines the length of generated OTPs.
	OTPLength = 6
	// OTPExpiryMinutes defines how long an OTP is valid.
	OTPExpiryMinutes = 5
)

// GenerateOTP creates a cryptographically secure 6-digit OTP.
func GenerateOTP() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to less secure method if crypto/rand fails
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// CreateOTPRecord creates a new OTP record for a phone number.
func CreateOTPRecord(phone string) *models.OTPRecord {
	return &models.OTPRecord{
		Phone:     phone,
		OTP:       GenerateOTP(),
		ExpiresAt: time.Now().Add(time.Duration(OTPExpiryMinutes) * time.Minute),
	}
}

// MockSendOTP simulates sending an OTP to a phone number.
// Replace with actual SMS service integration in production.
func MockSendOTP(phone, otp string) error {
	log.Printf("📱 [MOCK SMS] Sending OTP to %s: %s", phone, otp)
	log.Printf("📱 [MOCK SMS] OTP expires in %d minutes", OTPExpiryMinutes)
	// In production, integrate with SMS providers like:
	// - Twilio
	// - Africa's Talking
	// - Termii
	// - etc.
	return nil
}

// ValidateOTP checks if the provided OTP matches and is valid.
func ValidateOTP(record *models.OTPRecord, providedOTP string) bool {
	if record == nil {
		return false
	}
	if record.IsUsed {
		return false
	}
	if record.IsExpired() {
		return false
	}
	return record.OTP == providedOTP
}
