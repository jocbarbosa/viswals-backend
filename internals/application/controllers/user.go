package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jocbarbosa/viswals-backend/internals/core/dto/filters"
	"github.com/jocbarbosa/viswals-backend/internals/core/port"
)

type UserController struct {
	userService port.UserService
	logger      port.Logger
	cache       port.Cache
}

// NewUserController creates a new UserController
func NewUserController(svc port.UserService, logger port.Logger, cache port.Cache) *UserController {
	return &UserController{
		userService: svc,
		logger:      logger,
		cache:       cache,
	}
}

// GetUsers is the handler for the GET /users route
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("handling GET /users request")

	ctx := r.Context()

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	email := r.URL.Query().Get("email")

	filters := filters.UserFilter{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	users, err := c.userService.GetUsers(ctx, filters)
	if err != nil {
		c.logger.Error("error finding all users", err)
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		c.logger.Error("error encoding users response", err)
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
}
