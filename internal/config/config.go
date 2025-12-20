package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

var AppConfig *Config

func LoadConfig() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	AppConfig = &Config{
		JWTSecret:     jwtSecret,
		JWTExpiration: 24 * time.Hour, // Token expires in 24 hours
	}

	log.Println("Configuration loaded successfully")
}
