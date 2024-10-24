package repository

import (
	"errors"

	"github.com/jocbarbosa/viswals-backend/internals/core/dto/filters"
	"github.com/jocbarbosa/viswals-backend/internals/core/model"
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GormUserRepository
func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &GormUserRepository{db: db}
}

// FindByID finds a user by ID from repository
func (r *GormUserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users from repository using filters when available
func (r *GormUserRepository) FindAll(filters filters.UserFilter) ([]model.User, error) {
	var users []model.User
	query := r.db.Model(&model.User{})

	if filters.FirstName != "" {
		query = query.Where("first_name LIKE ?", "%"+filters.FirstName+"%")
	}
	if filters.LastName != "" {
		query = query.Where("last_name LIKE ?", "%"+filters.LastName+"%")
	}
	if filters.Email != "" {
		query = query.Where("email LIKE ?", "%"+filters.Email+"%")
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Create inserts a new user on repository
func (r *GormUserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update modifies an existing user on repository
func (r *GormUserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user by ID from repository
func (r *GormUserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
