package services

import (
	"context"
	"encoding/json"

	"github.com/jocbarbosa/viswals-backend/internals/core/model"
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
)

type UserService struct {
	logger    port.Logger
	userRepo  port.UserRepository
	cache     port.Cache
	messaging port.Messaging
}

// NewUserService creates a new UserService
func NewUserService(logger port.Logger, userRepo port.UserRepository, cache port.Cache, msg port.Messaging) *UserService {
	return &UserService{
		logger:    logger,
		userRepo:  userRepo,
		cache:     cache,
		messaging: msg,
	}
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Error("Error finding user by ID", err)
		return nil, err
	}
	return user, nil
}

// GetUsers returns all users
func (s *UserService) GetUsers() ([]model.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		s.logger.Error("Error finding all users", err)
		return nil, err
	}
	return users, nil
}

// StartConsuming starts consuming messages from messaging source
func (s *UserService) StartConsuming(ctx context.Context) {
	handler := func(msg port.Message) error {
		var user model.User

		err := json.Unmarshal(msg.Value, &user)
		if err != nil {
			s.logger.Error("Failed to unmarshal user", err)
			return err
		}

		if err := s.userRepo.Create(&user); err != nil {
			s.logger.Error("Failed to store user in PostgreSQL", err)
			return err
		}
		s.logger.Info("User stored in PostgreSQL", user.ID)

		userKey := "user:" + string(rune(user.ID))
		userData, err := json.Marshal(user)
		if err != nil {
			s.logger.Error("Failed to marshal user for Redis", err)
			return err
		}

		err = s.cache.Set(ctx, userKey, userData, 3600)
		if err != nil {
			s.logger.Error("Failed to store user in Redis", err)
			return err
		}
		s.logger.Info("User cached in Redis", user.ID)

		err = msg.AckFunc()
		if err != nil {
			s.logger.Error("Failed to acknowledge RabbitMQ message", err)
			return err
		}

		return nil
	}

	if err := s.messaging.Consume(handler); err != nil {
		s.logger.Error("Failed to start consuming RabbitMQ messages", err)
	}
}
