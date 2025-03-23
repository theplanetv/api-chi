package config

import "os"

var (
	WEB_URL string
)

func LoadWebConfig() {
	WEB_URL = os.Getenv("WEB_URL")
}
