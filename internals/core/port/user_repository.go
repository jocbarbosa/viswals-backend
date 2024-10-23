package port

import "github.com/jocbarbosa/viswals-backend/internals/core/model"

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	Delete(id uint) error
}
