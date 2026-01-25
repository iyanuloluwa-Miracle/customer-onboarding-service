package models

// OnboardRequest represents the request body for customer onboarding.
// @Description Request payload for customer onboarding
type OnboardRequest struct {
	Phone    string `json:"phone" binding:"required" example:"+2348012345678"`
	Email    string `json:"email" binding:"required,email" example:"customer@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"securepassword123"`
	State    string `json:"state" binding:"required" example:"Lagos"`
	LGA      string `json:"lga" binding:"required" example:"Ikeja"`
}

// VerifyOTPRequest represents the request body for OTP verification.
// @Description Request payload for OTP verification
type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required" example:"+2348012345678"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// ResendOTPRequest represents the request body for resending OTP.
// @Description Request payload for resending OTP
type ResendOTPRequest struct {
	Phone string `json:"phone" binding:"required" example:"+2348012345678"`
}

// APIResponse represents a standard API response.
// @Description Standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:""`
}

// CustomerResponse represents customer data in API responses.
// @Description Customer data in API responses
type CustomerResponse struct {
	ID              string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Phone           string `json:"phone" example:"+2348012345678"`
	Email           string `json:"email" example:"customer@example.com"`
	State           string `json:"state" example:"Lagos"`
	LGA             string `json:"lga" example:"Ikeja"`
	IsPhoneVerified bool   `json:"isPhoneVerified" example:"true"`
	CreatedAt       string `json:"createdAt" example:"2026-01-25T22:00:00Z"`
}

// Bank represents a bank from the external API.
// @Description Bank information from external API
type Bank struct {
	Code string `json:"code" example:"044"`
	Name string `json:"name" example:"Access Bank"`
}
