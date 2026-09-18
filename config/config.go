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
		log.Println("no .env file found, relying on system environment variables")
	}
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		UserA:           getEnv("USER_A", "alice"),
		UserB:           getEnv("USER_B", "bob"),
		AllowAllOrigins: getEnv("ALLOW_ALL_ORIGINS", "false") == "true",
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
