package config

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ServiceName       string `mapstructure:"SERVICE_NAME"`
	ServiceDomain     string `mapstructure:"SERVICE_DOMAIN"`
	ApiPrefix         string `mapstructure:"API_PREFIX"`
	HttpPort          string `mapstructure:"HTTP_PORT"`
	RequestTimeout    string `mapstructure:"REQUEST_TIMEOUT"`
	DatabaseURI       string `mapstructure:"DATABASE_URI"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	RedisAddress      string `mapstructure:"REDIS_ADDRESS"`
	MailerHost        string `mapstructure:"MAILER_HOST"`
	MailerPort        string `mapstructure:"MAILER_PORT"`
	MailerSEC         string `mapstructure:"MAILER_SEC"`
	MailerUserName    string `mapstructure:"MAILER_USERNAME"`
	MailerPassword    string `mapstructure:"MAILER_PASSWORD"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

func SetupEnvironment() (AppConfig, error) {
	env := getEnvWithDefault("APP_ENV", "production")
	log.Printf("[DEBUG] application running in: %s", env)

	if env == "dev" {
		return loadDevelopmentConfig()
	}

	return loadProductionConfig()
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadDevelopmentConfig() (AppConfig, error) {
	var appConfig AppConfig

	log.Println("[DEBUG] loading development environment from .env file")

	viper.SetConfigFile("example.env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return AppConfig{}, errors.New("could not read .env file")
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		return AppConfig{}, errors.New("could not unmarshal .env file")
	}

	requiredVars := map[string]*string{
		"SERVICE_NAME":        &appConfig.ServiceName,
		"SERVICE_DOMAIN":      &appConfig.ServiceDomain,
		"API_PREFIX":          &appConfig.ApiPrefix,
		"HTTP_PORT":           &appConfig.HttpPort,
		"REQUEST_TIMEOUT":     &appConfig.RequestTimeout,
		"DATABASE_URI":        &appConfig.DatabaseURI,
		"MIGRATION_URL":       &appConfig.MigrationURL,
		"REDIS_ADDRESS":       &appConfig.RedisAddress,
		"MAILER_HOST":         &appConfig.MailerHost,
		"MAILER_PORT":         &appConfig.MailerPort,
		"MAILER_SEC":          &appConfig.MailerSEC,
		"MAILER_USERNAME":     &appConfig.MailerUserName,
		"MAILER_PASSWORD":     &appConfig.MailerPassword,
		"TOKEN_SYMMETRIC_KEY": &appConfig.TokenSymmetricKey,
	}

	for key, value := range requiredVars {
		if *value == "" {
			return AppConfig{}, errors.New(key + " environment variable not found in .env file")
		}
	}

	if len(appConfig.TokenSymmetricKey) != 32 {
		return AppConfig{}, errors.New("TOKEN_SYMMETRIC_KEY must be exactly 32 characters long")
	}

	log.Println("[DEBUG] development environment loaded successfully!")
	return appConfig, nil
}

func loadProductionConfig() (AppConfig, error) {
	var appConfig AppConfig

	requiredVars := map[string]*string{
		"SERVICE_NAME":        &appConfig.ServiceName,
		"SERVICE_DOMAIN":      &appConfig.ServiceDomain,
		"API_PREFIX":          &appConfig.ApiPrefix,
		"HTTP_PORT":           &appConfig.HttpPort,
		"REQUEST_TIMEOUT":     &appConfig.RequestTimeout,
		"DATABASE_URI":        &appConfig.DatabaseURI,
		"MIGRATION_URL":       &appConfig.MigrationURL,
		"REDIS_ADDRESS":       &appConfig.RedisAddress,
		"MAILER_HOST":         &appConfig.MailerHost,
		"MAILER_PORT":         &appConfig.MailerPort,
		"MAILER_SEC":          &appConfig.MailerSEC,
		"MAILER_USERNAME":     &appConfig.MailerUserName,
		"MAILER_PASSWORD":     &appConfig.MailerPassword,
		"TOKEN_SYMMETRIC_KEY": &appConfig.TokenSymmetricKey,
	}

	for key, value := range requiredVars {
		if *value = os.Getenv(key); *value == "" {
			return AppConfig{}, errors.New(key + " environment variable not found")
		}
	}

	if len(appConfig.TokenSymmetricKey) != 32 {
		return AppConfig{}, errors.New("TOKEN_SYMMETRIC_KEY must be exactly 32 characters long")
	}

	log.Println("[DEBUG] Production environment variables loaded successfully")
	return appConfig, nil
}
