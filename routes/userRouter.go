package routes

import (
	"github.com/gorilla/mux"
	controllers "github.com/snehadeep-wagh/go-backend/controllers"
	"github.com/snehadeep-wagh/go-backend/middleware"
)

func UserRoutes(r *mux.Router) {
	// r.Use(middleware.Authenticate)
	userSubRoutes := r.Methods("GET").Subrouter()
	userSubRoutes.Use(middleware.Authenticate)
	userSubRoutes.HandleFunc("/users", controllers.GetUser)
	userSubRoutes.HandleFunc("/users/{userId}", controllers.GetUserById)
}
