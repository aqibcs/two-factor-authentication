package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// PostgreSQL connection parameters
var (
	HttpPort = 3636

	DbUser = "myuser"
	DbPass = "my_password"
	DbName = "mydatabase"
	DbHost = "localhost"
	DbPort = "5432"
)

// init loads environment variables from a .env file and overrides default values.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	if port, err := strconv.Atoi(os.Getenv("HTTP_PORT")); err != nil {
		HttpPort = port
	}

	if x := os.Getenv("DB_USER"); x != "" {
		DbUser = x
	}

	if x := os.Getenv("DB_PASSWORD"); x != "" {
		DbPass = x
	}

	if x := os.Getenv("DB_NAME"); x != "" {
		DbName = x
	}

	if x := os.Getenv("DB_HOST"); x != "" {
		DbHost = x
	}

	if x := os.Getenv("DB_PORT"); x != "" {
		DbPort = x
	}
}
