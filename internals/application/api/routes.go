package api

import (
	"net/http"

	"github.com/jocbarbosa/viswals-backend/internals/application/controllers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter sets up the API routes
func NewRouter(userController *controllers.UserController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", userController.GetUsers)

	mux.HandleFunc("/docs/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yaml")
	})

	mux.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.yaml"),
	))

	return mux
}
