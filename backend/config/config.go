package config

import (
	"log"

	"github.com/spf13/viper"
)

// DBConfig holds all database-related configuration.
type DBConfig struct {
	Driver   string `mapstructure:"DB_DRIVER"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

// Config is the main configuration struct for your application.
type Config struct {
	DB            DBConfig `mapstructure:",squash"`
	AppEnv        string   `mapstructure:"APP_ENV"`
	ServerAddress string   `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

