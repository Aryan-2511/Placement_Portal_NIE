package utils

import (
	"net/http"
)
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins (or specify a domain)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Specify allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Specify allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Role")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight (OPTIONS) requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}