package routes

import(
	"github.com/gorilla/mux"
	controllers "github.com/snehadeep-wagh/go-backend/controllers"
	"github.com/snehadeep-wagh/go-backend/middleware"
)

func UserRoutes(r *mux.Router){
	r.Use(middleware.Authenticate)
	r.HandleFunc("/users", controllers.GetUser).Methods("GET")
	r.HandleFunc("/users/{userId}", controllers.GetUserById).Methods("GET")
}