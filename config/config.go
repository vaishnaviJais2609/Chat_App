package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	AllowAllOrigins bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables and defaults")
	}

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		AllowAllOrigins: getEnv("ALLOW_ALL_ORIGINS", "true") == "true",
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
