package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	routes "github.com/snehadeep-wagh/go-backend/routes"
)

func main() {
	// Get the port from the env file
	port := os.Getenv("PORT")
	if port == "" { // check if the port is empty
		port = "8000" // assign port with default :8000 port
	}

	// create new instance of router
	r := mux.NewRouter()
	
	routes.AuthRoutes(r)
	routes.UserRoutes(r)

	fmt.Print("Running server on port: " + port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
