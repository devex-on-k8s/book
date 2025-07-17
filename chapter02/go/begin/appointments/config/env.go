package config

import (
	"fmt"
	"os"
)

var (
	Version            = getEnv("VERSION", "0.0.1-SNAPSHOT")
	Source             = getEnv("SOURCE", "https://github.com/")
	AppPort            = getEnv("APP_PORT", "8081")
	PostgresqlHost     = getEnv("POSTGRES_HOST", "localhost")
	PostgresqlPort     = getEnv("POSTGRES_PORT", "5432")
	PostgresqlDatabase = getEnv("POSTGRES_DB", "postgres")
	PostgresqlUsername = getEnv("POSTGRES_USERNAME", "postgres")
	PostgresqlPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	DB                 string
)

func getEnv(name string, fallback string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	if fallback != "" {
		return fallback
	}

	// Make sure the environment variables are set.
	panic(fmt.Sprintf(`Environment variable not found :: %v`, name))
}
