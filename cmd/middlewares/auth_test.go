package middlewares

import (
	"api-chi/cmd/config"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a token for testing
func createTestToken(username string, expiration time.Duration, secretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func TestCheckLogin(t *testing.T) {
	//authService := &services.AuthService{}
	authMiddleware := AuthMiddleware{} //service: authService}

	// Mock a next handler to test middleware chain behavior
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Authenticated!")
	})

	t.Run("Authorized request - valid token", func(t *testing.T) {
		// Create a valid token
		secretKey := []byte(config.AUTH_SECRET_KEY)
		tokenString, err := createTestToken("admin", 1*time.Hour, secretKey)
		assert.NoError(t, err)

		// Create a request with the token in a cookie
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "auth-token", Value: tokenString})
		rr := httptest.NewRecorder()

		// Run the middleware
		handler := authMiddleware.CheckLogin(nextHandler)
		handler.ServeHTTP(rr, req)

		// Assert that the response is as expected
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Authenticated!")
	})

	t.Run("Unauthorized request - missing token", func(t *testing.T) {
		// Create a request without the Authorization cookie
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		// Run the middleware
		handler := authMiddleware.CheckLogin(nextHandler)
		handler.ServeHTTP(rr, req)

		// Assert that the response is Unauthorized
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized: Missing token")
	})

	t.Run("Unauthorized request - expired token", func(t *testing.T) {
		// Create an expired token
		secretKey := []byte(config.AUTH_SECRET_KEY)
		tokenString, err := createTestToken("admin", -1*time.Hour, secretKey) // Token expired one hour ago
		assert.NoError(t, err)

		// Create a request with the expired token in a cookie
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "auth-token", Value: tokenString})
		rr := httptest.NewRecorder()

		// Run the middleware
		handler := authMiddleware.CheckLogin(nextHandler)
		handler.ServeHTTP(rr, req)

		// Assert that the response is Unauthorized
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized")
	})
}
