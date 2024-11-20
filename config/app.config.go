package config

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/viper"
)

// AppConfig holds the configuration values for the application.
// Each field is mapped to an environment variable using the `mapstructure` tag.
type AppConfig struct {
	ServiceName       string `mapstructure:"SERVICE_NAME"`
	ServiceDomain     string `mapstructure:"SERVICE_DOMAIN"`
	ApiPrefix         string `mapstructure:"API_PREFIX"`
	HttpPort          string `mapstructure:"HTTP_PORT"`
	RequestTimeout    string `mapstructure:"REQUEST_TIMEOUT"`
	DatabaseURI       string `mapstructure:"DATABASE_URI"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

// SetupEnvironment reads and sets up the application environment configuration.
// It checks if the `APP_ENV` variable is set and reads from `.env` if in development mode.
// Returns an AppConfig struct and an error if the setup fails.
func SetupEnvironment() (AppConfig, error) {
	var appConfig AppConfig

	// Get the environment mode (e.g., "dev" or "production").
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production" // Default to production mode if not set.
	}
	log.Printf("[DEBUG] Environment set to: %s", env)

	// Check if the environment is set to development mode.
	if env == "dev" {
		log.Println("[DEBUG] Loading development environment from .env file")

		// In development mode, configure Viper to read from a `.env` file.
		viper.AddConfigPath(".")
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		// Attempt to read the configuration from the `.env` file.
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("[ERROR] Error reading .env file: %v", err)
			return AppConfig{}, errors.New("could not read environment variables")
		}
		log.Println("[DEBUG] .env file loaded successfully")

		// Unmarshal the config values into the `appConfig` struct.
		if err := viper.Unmarshal(&appConfig); err != nil {
			log.Printf("[ERROR] Error unmarshalling .env data: %v", err)
			return AppConfig{}, errors.New("could not unmarshal environment variables")
		}
	} else {
		// In production mode, read essential environment variables directly.
		log.Println("[DEBUG] Loading production environment variables")

		// Read `HTTP_PORT` and check if it is set.
		httpPort := os.Getenv("HTTP_PORT")
		if len(httpPort) < 1 {
			log.Println("[ERROR] HTTP_PORT environment variable not found")
			return AppConfig{}, errors.New("HTTP_PORT environment variable not found")
		}

		// Read `DATABASE_URI` and check if it is set.
		databaseURI := os.Getenv("DATABASE_URI")
		if len(databaseURI) < 1 {
			log.Println("[ERROR] DATABASE_URI environment variable not found")
			return AppConfig{}, errors.New("DATABASE_URI environment variable not found")
		}

		// Read `TOKEN_SYMMETRIC_KEY` and validate its length.
		tokenSymmetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")
		if len(tokenSymmetricKey) != 32 {
			log.Println("[ERROR] TOKEN_SYMMETRIC_KEY is not exactly 32 characters")
			return AppConfig{}, errors.New("TOKEN_SYMMETRIC_KEY must be exactly 32 characters long")
		}

		// Set the read values into the `appConfig` struct.
		appConfig.HttpPort = httpPort
		appConfig.DatabaseURI = databaseURI
		appConfig.TokenSymmetricKey = tokenSymmetricKey

		log.Println("[DEBUG] Production environment variables loaded successfully")
	}

	// Return the populated AppConfig struct and nil (no error).
	return appConfig, nil
}
