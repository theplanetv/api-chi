package services

import (
	"api-chi/cmd/config"
	"api-chi/cmd/models"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (s *AuthService) New() {
	config.LoadAuthConfig()
}

func (s *AuthService) GenerateToken(input *models.Auth) (string, error) {
	// Set token expiration time (e.g., 1 hour)
	expirationTime := time.Now().Add(1 * time.Hour).Unix()

	// Create JWT claims, which includes the user ID and expiration time
	claims := jwt.MapClaims{
		"username": input.Username,
		"exp":      expirationTime,
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Retrieve the secret key from config
	secretKey := []byte(config.AUTH_SECRET_KEY)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (bool, error) {
	// Retrieve the secret key from config
	secretKey := []byte(config.AUTH_SECRET_KEY)

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is as expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Check if token is valid and not expired
	if err != nil || !token.Valid {
		return false, fmt.Errorf("invalid or expired token")
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verify expiration time
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, fmt.Errorf("token expiration not found")
		}

		// Check if token is expired
		if int64(exp) < time.Now().Unix() {
			return false, fmt.Errorf("token has expired")
		}

		// Extract username
		_, ok = claims["username"].(string)
		if !ok {
			return false, fmt.Errorf("username claim not found")
		}

		// Return the authenticated user
		return true, nil
	}

	return false, fmt.Errorf("invalid token claims")
}

func (s *AuthService) Login(input *models.Auth) error {
	if config.AUTH_BCRYPT_COST == "" {
		config.AUTH_BCRYPT_COST = "10"
	}
	bcryptCost, err := strconv.Atoi(config.AUTH_BCRYPT_COST)
	if err != nil {
		return fmt.Errorf("can't set bcrypt default cost")
	}

	// Check username
	if input.Username != config.AUTH_USERNAME {
		return fmt.Errorf("username is incorrect")
	}

	// Check password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err)
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(config.AUTH_PASSWORD))
	if err != nil {
		return fmt.Errorf("password is incorrect")
	}

	return nil
}
