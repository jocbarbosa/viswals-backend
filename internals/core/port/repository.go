package port

import (
	"github.com/jocbarbosa/viswals-backend/internals/core/dto/filters"
	"github.com/jocbarbosa/viswals-backend/internals/core/model"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	Upsert(user *model.User) (*model.User, error)
	Update(user *model.User) error
	FindAll(filters filters.UserFilter) ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	Delete(id uint) error
}
