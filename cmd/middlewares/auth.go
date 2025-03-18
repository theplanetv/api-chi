package middlewares

import (
	"api-chi/cmd/services"
	"fmt"
	"net/http"
)

type AuthMiddleware struct {
	service services.AuthService
}

func (m *AuthMiddleware) CheckLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the "Authorization" cookie
		cookie, err := r.Cookie("auth-token")
		if err != nil {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		// Validate the token using the service
		isValid, err := m.service.ValidateToken(tokenString)
		if err != nil || !isValid {
			// If token validation fails, respond with unauthorized status
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if the token is valid
		next.ServeHTTP(w, r)
	})
}
