package infra

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		Msg.Errorf("Error loading .env file: %s", err.Error())
	}
}

type Config struct {
	Port               string
	DatabaseConnection string
}

func GetConfig() *Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	conn := os.Getenv("DB_CONNECTION")

	if conn == "" {
		Msg.Errorf("Database connection is not set, using default value %s", conn)
	}

	return &Config{
		Port:               port,
		DatabaseConnection: conn,
	}
}
