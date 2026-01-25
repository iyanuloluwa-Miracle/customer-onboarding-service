package repositories

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gin-quickstart/internal/models"
)

// OTPRepository implements OTPRepositoryInterface.
type OTPRepository struct {
	db      *gorm.DB
	records []models.OTPRecord // In-memory fallback
	mu      sync.RWMutex
}

// NewOTPRepository creates a new OTPRepository.
func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{
		db:      db,
		records: make([]models.OTPRecord, 0),
	}
}

// Create adds a new OTP record.
func (r *OTPRepository) Create(otp *models.OTPRecord) error {
	if r.db != nil {
		return r.db.Create(otp).Error
	}

	// In-memory fallback
	r.mu.Lock()
	defer r.mu.Unlock()

	if otp.ID == uuid.Nil {
		otp.ID = uuid.New()
	}
	if otp.CreatedAt.IsZero() {
		otp.CreatedAt = time.Now()
	}
	r.records = append(r.records, *otp)
	return nil
}

// FindLatestByPhone finds the latest unused OTP for a phone number.
func (r *OTPRepository) FindLatestByPhone(phone string) (*models.OTPRecord, error) {
	if r.db != nil {
		var otp models.OTPRecord
		err := r.db.Where("phone = ? AND is_used = ? AND expires_at > ?", phone, false, time.Now()).
			Order("created_at DESC").
			First(&otp).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return &otp, err
	}

	// In-memory fallback
	r.mu.RLock()
	defer r.mu.RUnlock()

	var latest *models.OTPRecord
	for i := range r.records {
		rec := &r.records[i]
		if rec.Phone == phone && !rec.IsUsed && !rec.IsExpired() {
			if latest == nil || rec.CreatedAt.After(latest.CreatedAt) {
				latest = rec
			}
		}
	}
	return latest, nil
}

// MarkAsUsed marks an OTP as used.
func (r *OTPRepository) MarkAsUsed(id uuid.UUID) error {
	if r.db != nil {
		return r.db.Model(&models.OTPRecord{}).Where("id = ?", id).Update("is_used", true).Error
	}

	// In-memory fallback
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.records {
		if r.records[i].ID == id {
			r.records[i].IsUsed = true
			return nil
		}
	}
	return nil
}

// DeleteExpired removes expired OTP records.
func (r *OTPRepository) DeleteExpired() error {
	if r.db != nil {
		return r.db.Where("expires_at < ?", time.Now()).Delete(&models.OTPRecord{}).Error
	}

	// In-memory fallback
	r.mu.Lock()
	defer r.mu.Unlock()

	var valid []models.OTPRecord
	for _, rec := range r.records {
		if !rec.IsExpired() {
			valid = append(valid, rec)
		}
	}
	r.records = valid
	return nil
}
