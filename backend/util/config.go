package util

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	AppEnv        string `mapstructure:"APP_ENV"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig() (Config, error) {
	// Enable environment variables
	viper.AutomaticEnv()

	// Explicitly bind environment variables for database config
	viper.BindEnv("Database.Host", "POSTGRES_HOST")
	viper.BindEnv("Database.Port", "POSTGRES_PORT")
	viper.BindEnv("Database.User", "POSTGRES_USER")
	viper.BindEnv("Database.Password", "POSTGRES_PASSWORD")
	viper.BindEnv("Database.DbName", "POSTGRES_DB")

	// Try to load .env file if it exists (for development)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only log if error is not "config file not found"
			log.Printf("Error reading config file: %v", err)
		}
	}

	// Set default values
	viper.SetDefault("APP_PORT", "8081")
	viper.SetDefault("APP_ENV", "development")

	var config Config

	// Manually set database config from environment variables if they exist
	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		config.DBHost = host
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		config.DBPort = port
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		config.DBUser = user
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		config.DBPassword = password
	}
	if dbName := os.Getenv("POSTGRES_DB"); dbName != "" {
		config.DBName = dbName
	}

	// Unmarshal the rest of the config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
