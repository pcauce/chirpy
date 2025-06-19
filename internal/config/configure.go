package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const Port = "8080"

type ApiConfig struct {
	Platform      string
	JWTSecret     string
	TokenDuration map[string]time.Duration
	PolkaKey      string
}

var api ApiConfig

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api = ApiConfig{
		Platform:  os.Getenv("PLATFORM"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		TokenDuration: map[string]time.Duration{
			"access":  time.Hour,
			"refresh": time.Hour * 24 * 60,
		},
		PolkaKey: os.Getenv("POLKA_KEY"),
	}
}

func APIConfig() *ApiConfig {
	return &api
}
