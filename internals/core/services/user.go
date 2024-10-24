package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jocbarbosa/viswals-backend/internals/config"
	"github.com/jocbarbosa/viswals-backend/internals/core/dto/filters"
	"github.com/jocbarbosa/viswals-backend/internals/core/model"
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
	"github.com/jocbarbosa/viswals-backend/internals/utils"
)

type UserService struct {
	logger    port.Logger
	userRepo  port.UserRepository
	cache     port.Cache
	messaging port.Messaging
	config    config.Config
}

// NewUserService creates a new UserService
func NewUserService(logger port.Logger, userRepo port.UserRepository, cache port.Cache, msg port.Messaging) *UserService {
	return &UserService{
		logger:    logger,
		userRepo:  userRepo,
		cache:     cache,
		messaging: msg,
		config:    config.NewConfig(),
	}
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Error("error finding user by ID", err)
		return nil, err
	}

	user.Email, err = utils.Decrypt(user.Email, s.config.EncryptionKey)
	if err != nil {
		s.logger.Error("failed to decrypt user email", err)
		return nil, err
	}

	return user, nil
}

// GetUsers returns all users
func (s *UserService) GetUsers(ctx context.Context, filters filters.UserFilter) ([]model.User, error) {
	if filters.Email != "" {
		encryptedEmail, err := utils.Encrypt(filters.Email, s.config.EncryptionKey)
		if err != nil {
			s.logger.Error("failed to encrypt email filter", err)
			return nil, err
		}
		filters.Email = encryptedEmail

		userKey := fmt.Sprintf("user:%s", encryptedEmail)
		cachedUser, err := s.cache.Get(ctx, userKey)
		if err == nil && cachedUser != nil {
			cachedUserStr, ok := cachedUser.(string)
			if ok {
				var user model.User
				if err := json.Unmarshal([]byte(cachedUserStr), &user); err == nil {
					s.logger.Info("retrieved user from cache", user.Email)
					return []model.User{user}, nil
				}
			}
		}
	}

	users, err := s.userRepo.FindAll(filters)
	if err != nil {
		s.logger.Error("error finding all users", err)
		return nil, err
	}

	for i := range users {
		users[i].Email, err = utils.Decrypt(users[i].Email, s.config.EncryptionKey)
		if err != nil {
			s.logger.Error("failed to decrypt user email", err)
		}
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

		user.Email, err = utils.Encrypt(user.Email, s.config.EncryptionKey)
		if err != nil {
			s.logger.Error("failed to encrypt user email", err)
			return err
		}

		createdUser, err := s.userRepo.Upsert(&user)
		if err != nil {
			s.logger.Error("failed to store user in repository", err)
			return err
		}
		s.logger.Info("user stored in repository", createdUser.ID)

		userKey := "user:" + user.Email

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
		s.logger.Info("user cached successfully", createdUser.ID)

		if msg.AckFunc != nil {
			err = msg.AckFunc()
			if err != nil {
				s.logger.Error("failed to acknowledge message", err)
				return err
			}
		}

		return nil
	}

	err := s.messaging.Consume(handler)
	if err != nil {
		s.logger.Error("failed to start consuming messages", err)
	}
}
