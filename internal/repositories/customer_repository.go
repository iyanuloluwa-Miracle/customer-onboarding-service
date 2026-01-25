package repositories

import (
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gin-quickstart/internal/models"
)

// CustomerRepository implements CustomerRepositoryInterface.
// It supports both GORM (MS SQL) and in-memory storage.
type CustomerRepository struct {
	db        *gorm.DB
	customers []models.Customer // In-memory fallback
	mu        sync.RWMutex
}

// NewCustomerRepository creates a new CustomerRepository.
// If db is nil, it uses in-memory storage.
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		db:        db,
		customers: make([]models.Customer, 0),
	}
}

// Create adds a new customer to the repository.
func (r *CustomerRepository) Create(customer *models.Customer) error {
	if r.db != nil {
		return r.db.Create(customer).Error
	}

	// In-memory fallback
	r.mu.Lock()
	defer r.mu.Unlock()

	if customer.ID == uuid.Nil {
		customer.ID = uuid.New()
	}
	r.customers = append(r.customers, *customer)
	return nil
}

// GetAll retrieves all customers from the repository.
func (r *CustomerRepository) GetAll() ([]models.Customer, error) {
	if r.db != nil {
		var customers []models.Customer
		err := r.db.Find(&customers).Error
		return customers, err
	}

	// In-memory fallback
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.customers, nil
}

// FindByID finds a customer by ID.
func (r *CustomerRepository) FindByID(id uuid.UUID) (*models.Customer, error) {
	if r.db != nil {
		var customer models.Customer
		err := r.db.Where("id = ?", id).First(&customer).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return &customer, err
	}

	// In-memory fallback
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.customers {
		if r.customers[i].ID == id {
			return &r.customers[i], nil
		}
	}
	return nil, nil
}

// FindByPhone finds a customer by phone number.
func (r *CustomerRepository) FindByPhone(phone string) (*models.Customer, error) {
	if r.db != nil {
		var customer models.Customer
		err := r.db.Where("phone = ?", phone).First(&customer).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return &customer, err
	}

	// In-memory fallback
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.customers {
		if r.customers[i].Phone == phone {
			return &r.customers[i], nil
		}
	}
	return nil, nil
}

// FindByEmail finds a customer by email address.
func (r *CustomerRepository) FindByEmail(email string) (*models.Customer, error) {
	if r.db != nil {
		var customer models.Customer
		err := r.db.Where("email = ?", email).First(&customer).Error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return &customer, err
	}

	// In-memory fallback
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.customers {
		if r.customers[i].Email == email {
			return &r.customers[i], nil
		}
	}
	return nil, nil
}

// Update updates an existing customer.
func (r *CustomerRepository) Update(customer *models.Customer) error {
	if r.db != nil {
		return r.db.Save(customer).Error
	}

	// In-memory fallback
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.customers {
		if r.customers[i].ID == customer.ID {
			r.customers[i] = *customer
			return nil
		}
	}
	return nil
}

// Delete soft-deletes a customer by ID.
func (r *CustomerRepository) Delete(id uuid.UUID) error {
	if r.db != nil {
		return r.db.Delete(&models.Customer{}, id).Error
	}

	// In-memory fallback
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.customers {
		if r.customers[i].ID == id {
			r.customers = append(r.customers[:i], r.customers[i+1:]...)
			return nil
		}
	}
	return nil
}
