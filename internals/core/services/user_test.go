package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jocbarbosa/viswals-backend/internals/core/dto/filters"
	"github.com/jocbarbosa/viswals-backend/internals/core/model"
	"github.com/jocbarbosa/viswals-backend/internals/core/services"
	"github.com/jocbarbosa/viswals-backend/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetUsers(t *testing.T) {
	mockLogger := &mocks.Logger{}
	mockRepo := &mocks.UserRepository{}
	mockCache := &mocks.Cache{}
	mockMessaging := &mocks.Messaging{}

	userService := services.NewUserService(mockLogger, mockRepo, mockCache, mockMessaging)

	tests := []struct {
		name    string
		filter  filters.UserFilter
		setup   func()
		want    []model.User
		wantErr bool
	}{
		{
			name:   "get user from cache",
			filter: filters.UserFilter{Email: "any@email.com"},
			setup: func() {
				user := &model.User{ID: 1, FirstName: "any", LastName: "any", Email: "any@email.com"}
				userData, _ := json.Marshal(user)
				mockCache.On("Get", mock.Anything, "user:any@email.com").Return(string(userData), nil).Once()
			},
			want:    []model.User{{ID: 1, FirstName: "any", LastName: "any", Email: "any@email.com"}},
			wantErr: false,
		},
		{
			name:   "get user from repository",
			filter: filters.UserFilter{},
			setup: func() {
				mockRepo.On("FindAll", filters.UserFilter{}).Return([]model.User{
					{ID: 1, FirstName: "any", LastName: "any", Email: "any@email.com"},
				}, nil).Once()
			},
			want:    []model.User{{ID: 1, FirstName: "any", LastName: "any", Email: "any@email.com"}},
			wantErr: false,
		},
		{
			name:   "get repository error",
			filter: filters.UserFilter{},
			setup: func() {
				mockRepo.On("FindAll", filters.UserFilter{}).Return(nil, errors.New("repo error")).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := userService.GetUsers(context.Background(), tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}
