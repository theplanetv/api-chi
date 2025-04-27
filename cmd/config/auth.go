package config

import "os"

var (
	// Auth config
	AUTH_USERNAME    string
	AUTH_PASSWORD    string
	AUTH_SECRET_KEY  string
	AUTH_BCRYPT_COST string
)

func LoadAuthConfig() {
	AUTH_USERNAME = os.Getenv("API_CHI_AUTH_USERNAME")
	AUTH_PASSWORD = os.Getenv("API_CHI_AUTH_PASSWORD")
	AUTH_SECRET_KEY = os.Getenv("API_CHI_AUTH_SECRET_KEY")
	AUTH_BCRYPT_COST = os.Getenv("API_CHI_AUTH_BCRYPT_COST")
}
