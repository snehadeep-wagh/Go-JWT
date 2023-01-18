package middleware

import (
	"net/http"

	"github.com/snehadeep-wagh/go-backend/helpers"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "No authentication token provided by user!"+token, http.StatusInternalServerError)
			return
		}

		_, err := helpers.ValidateToken(r, token)
		if err != "" {
			http.Error(w, err, http.StatusInternalServerError)
			return
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
