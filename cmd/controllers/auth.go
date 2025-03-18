package controllers

import (
	"api-chi/cmd/models"
	"api-chi/cmd/services"
	"api-chi/internal/message"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type AuthController struct {
	service services.AuthService
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// Get JSON from user input
	input := models.Auth{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, message.Response{
			Message: message.LOGIN_FAILED,
			Data:    nil,
		})
		return
	}

	// Execute Count and return if failed or success
	err := c.service.Login(&input)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, message.Response{
			Message: message.LOGIN_FAILED,
			Data:    nil,
		})
		return
	}

	// Generate JWT token after successful login
	token, err := c.service.GenerateToken(&input)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, message.Response{
			Message: message.LOGIN_FAILED,
			Data:    nil,
		})
		return
	}

	// Set auth-token as a cookie with the generated JWT token
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour), // Set a 24-hour expiration
		HttpOnly: true,                          // Prevents JavaScript access to the cookie
		Secure:   false,                         // Ensures the cookie is only sent over HTTPS
		SameSite: http.SameSiteStrictMode,       // Prevents CSRF by restricting cross-site cookie access
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, message.Response{
		Message: message.LOGIN_SUCCESS,
		Data:    nil,
	})
}

func (c *AuthController) Verify(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth-token")
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, message.Response{
			Message: message.AUTH_FAILED,
			Data:    nil,
		})
		return
	}

	if cookie.Value == "" {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, message.Response{
			Message: message.AUTH_FAILED,
			Data:    nil,
		})
		return
	}

	isAuth, err := c.service.ValidateToken(cookie.Value)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, message.Response{
			Message: message.AUTH_FAILED,
			Data:    nil,
		})
		return
	}

	if !isAuth {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, message.Response{
			Message: message.AUTH_FAILED,
			Data:    nil,
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, message.Response{
		Message: message.AUTH_SUCCESS,
		Data:    nil,
	})
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Invalidate the auth-token cookie by setting an expired date
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Set expiration to a time in the past
		HttpOnly: true,
		Secure:   false, // Use Secure flag if using HTTPS
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, message.Response{
		Message: message.LOGOUT_SUCCESS,
		Data:    nil,
	})
}
