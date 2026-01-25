package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Customer represents an onboarded customer in the system.
// @Description Customer entity with personal information and verification status
type Customer struct {
	ID              uuid.UUID      `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()" json:"id"`
	Phone           string         `gorm:"type:nvarchar(20);unique;not null" json:"phone"`
	Email           string         `gorm:"type:nvarchar(255);unique;not null" json:"email"`
	PasswordHash    string         `gorm:"type:nvarchar(255);not null" json:"-"`
	State           string         `gorm:"type:nvarchar(100);not null" json:"state"`
	LGA             string         `gorm:"type:nvarchar(100);not null" json:"lga"`
	IsPhoneVerified bool           `gorm:"default:false" json:"isPhoneVerified"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate sets UUID before creating a new customer.
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// OTPRecord stores OTP records for phone verification.
// @Description OTP record for phone verification
type OTPRecord struct {
	ID        uuid.UUID `gorm:"type:uniqueidentifier;primaryKey;default:NEWID()" json:"id"`
	Phone     string    `gorm:"type:nvarchar(20);not null;index" json:"phone"`
	OTP       string    `gorm:"type:nvarchar(10);not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	IsUsed    bool      `gorm:"default:false" json:"isUsed"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// BeforeCreate sets UUID before creating a new OTP record.
func (o *OTPRecord) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the OTP has expired.
func (o *OTPRecord) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}
