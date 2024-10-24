package port

import (
	"context"

	"github.com/jocbarbosa/viswals-backend/internals/core/dto/filters"
	"github.com/jocbarbosa/viswals-backend/internals/core/model"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUsers(ctx context.Context, filters filters.UserFilter) ([]model.User, error)
	StartConsuming(ctx context.Context)
}
