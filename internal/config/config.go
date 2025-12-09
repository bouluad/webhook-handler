package config

import (
	"log"
	"os"
)

// Config holds all application configuration settings.
type Config struct {
	// Server settings
	Port string

	// GitHub security
	GitHubSecret string

	// Azure Service Bus settings
	ServiceBusConnectionString string
	ServiceBusQueueName      string
}

// LoadConfig initializes the configuration from environment variables.
func LoadConfig() *Config {
	cfg := &Config{
		Port:                       getEnv("PORT", "8080"),
		GitHubSecret:               getEnvStrict("GITHUB_WEBHOOK_SECRET"),
		ServiceBusConnectionString: getEnvStrict("AZURE_SERVICE_BUS_CONN_STRING"),
		ServiceBusQueueName:        getEnvStrict("AZURE_SERVICE_BUS_QUEUE_NAME"),
	}
	log.Println("Configuration loaded successfully.")
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvStrict(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Environment variable %s is required and not set.", key)
	return "" // Should not reach here
}
