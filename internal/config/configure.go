package config

import (
	"os"
	"time"
)

const Port = "8080"

type apiConfig struct {
	Platform      string
	JWTSecret     string
	TokenDuration map[string]time.Duration
	PolkaKey      string
}

var API = apiConfig{
	Platform:  os.Getenv("PLATFORM"),
	JWTSecret: os.Getenv("JWT_SECRET"),
	TokenDuration: map[string]time.Duration{
		"access":  time.Hour,
		"refresh": time.Hour * 24 * 60,
	},
	PolkaKey: os.Getenv("POLKA_KEY"),
}
