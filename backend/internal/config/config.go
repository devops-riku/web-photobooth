package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	DOKey       string
	DOSecret    string
	DOEndpoint  string
	DORegion    string
	DOBucket    string
	GuestExpirationDays int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		DOKey:       os.Getenv("DO_SPACES_KEY"),
		DOSecret:    os.Getenv("DO_SPACES_SECRET"),
		DOEndpoint:  os.Getenv("DO_SPACES_ENDPOINT"),
		DORegion:    os.Getenv("DO_SPACES_REGION"),
		DOBucket:    os.Getenv("DO_SPACES_BUCKET"),
		GuestExpirationDays: 7, // Default to 7 days
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
