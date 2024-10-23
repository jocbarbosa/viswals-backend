package services

import (
	"context"
	"encoding/json"
	"strconv"

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
		s.logger.Error("error finding user by ID", err)
		return nil, err
	}
	return user, nil
}

// GetUsers returns all users
func (s *UserService) GetUsers() ([]model.User, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		s.logger.Error("error finding all users", err)
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
			s.logger.Error("failed to unmarshal user", err)
			return err
		}

		err = s.userRepo.Create(&user)
		if err != nil {
			s.logger.Error("failed to store user repository", err)
			return err
		}
		s.logger.Info("user stored in user repository", user.ID)

		userKey := "user:" + strconv.Itoa(user.ID)

		userData, err := json.Marshal(user)
		if err != nil {
			s.logger.Error("failed to marshal user for Redis", err)
			return err
		}

		err = s.cache.Set(ctx, userKey, userData, 3600) // keep for 1 hour
		if err != nil {
			s.logger.Error("failed to store user in Redis", err)
			return err
		}
		s.logger.Info("user cached with success", user.ID)

		err = msg.AckFunc()
		if err != nil {
			s.logger.Error("failed to acknowledge message", err)
			return err
		}

		return nil
	}

	err := s.messaging.Consume(handler)
	if err != nil {
		s.logger.Error("failed to start consuming messages", err)
	}
}
