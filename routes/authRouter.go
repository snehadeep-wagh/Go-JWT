package routes

import (
	"github.com/gorilla/mux"
	"github.com/snehadeep-wagh/go-backend/controllers"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("users/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("users/login", controllers.Login).Methods("POST")
}