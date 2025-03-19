package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PUBLIC_HOST string
	PORT        string

	DB_USER     string
	DB_PASSWORD string
	DB_ADDRESS  string
	DB_NAME     string

	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PUBLIC_HOST: getEnv("PUBLIC_HOST", "http://localhost"),
		PORT:        getEnv("PORT", "8080"),
		DB_USER:     getEnv("DB_USER", "root"),
		DB_PASSWORD: getEnv("DB_PASSWORD", "root"),
		DB_ADDRESS:  getEnv("DB_ADDRESS", "localhost:3306"),
		DB_NAME:     getEnv("DB_NAME", "ecommerce"),

		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
