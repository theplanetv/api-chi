package config

import "os"

var (
	// Api config
	API_PORT string
)

func LoadApiConfig() {
	API_PORT = os.Getenv("API_PORT")
}
