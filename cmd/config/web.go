package config

import "os"

var (
	// Web config
	WEB_HOST string
	WEB_PORT string
)

func LoadWebConfig() {
	WEB_HOST = os.Getenv("WEB_HOST")
	WEB_PORT = os.Getenv("WEB_PORT")
}
