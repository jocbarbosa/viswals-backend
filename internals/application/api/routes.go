package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/jocbarbosa/viswals-backend/internals/application/controllers"
)

// NewRouter sets up the API routes using chi
func NewRouter(userController *controllers.UserController) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/users", userController.GetUsers)

	return r
}
