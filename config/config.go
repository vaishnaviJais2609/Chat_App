package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	UserA           string
	UserB           string
	AllowAllOrigins bool
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:            getEnv("PORT", "8080"),
		UserA:           getEnv("USER_A", "alice"),
		UserB:           getEnv("USER_B", "bob"),
		AllowAllOrigins: getEnv("ALLOW_ALL_ORIGINS", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
