package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/snehadeep-wagh/go-backend/helpers"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		if token == "" {
			http.Error(w, "No authentication token provided!", http.StatusInternalServerError)
			return
		}

		_, err := helpers.ValidateToken(token)
		if err != "" {
			http.Error(w, err, http.StatusInternalServerError)
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
