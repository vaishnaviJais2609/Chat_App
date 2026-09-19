package config

import (
	"log"
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
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables and defaults")
	}

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		UserA:           getEnv("USER_A", "alice"),
		UserB:           getEnv("USER_B", "bob"),
		AllowAllOrigins: getEnv("ALLOW_ALL_ORIGINS", "false") == "true",
	}

	if cfg.UserA == cfg.UserB {
		log.Fatal("USER_A and USER_B must be different")
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
