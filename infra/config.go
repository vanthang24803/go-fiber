package infra

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseConnection string
	JWTSecret          string
	JWTRefreshSecret   string
}

var (
	config *Config
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %s\n", err)
	}
	config = loadConfig()
}

func loadConfig() *Config {
	return &Config{
		Port:               getEnv("PORT", "8000"),
		DatabaseConnection: getEnvMustExist("DB_CONNECTION"),
		JWTSecret:          getEnvMustExist("JWT_SECRET"),
		JWTRefreshSecret:   getEnvMustExist("JWT_REFRESH_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvMustExist(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("%s must be set", key))
	}
	return value
}

func GetConfig() *Config {
	return config
}
