package api

import (
	"net/http"

	"github.com/jocbarbosa/viswals-backend/internals/application/controllers"
)

// NewRouter sets up the API routes using chi
func NewRouter(userController *controllers.UserController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", userController.GetUsers)

	return mux
}
