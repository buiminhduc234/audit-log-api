package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort         int    `json:"server_port"`
	JWTSecretKey       string `json:"jwt_secret_key"`
	JWTExpirationHours int    `json:"jwt_expiration_hours"`
}

func Load() (*Config, error) {
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if serverPort == 0 {
		serverPort = 10000
	}

	jwtExpirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if jwtExpirationHours == 0 {
		jwtExpirationHours = 24
	}

	return &Config{
		ServerPort:         serverPort,
		JWTSecretKey:       os.Getenv("JWT_SECRET_KEY"),
		JWTExpirationHours: jwtExpirationHours,
	}, nil
}
