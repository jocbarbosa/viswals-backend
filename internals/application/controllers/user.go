package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jocbarbosa/viswals-backend/internals/core/port"
)

type UserController struct {
	userRepo port.UserRepository
	logger   port.Logger
	cache    port.Cache
}

// NewUserController creates a new UserController
func NewUserController(repo port.UserRepository, logger port.Logger, cache port.Cache) *UserController {
	return &UserController{
		userRepo: repo,
		logger:   logger,
		cache:    cache,
	}
}

// GetUsers is the handler for the GET /users route
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("handling GET /users request")

	users, err := c.userRepo.FindAll()
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
