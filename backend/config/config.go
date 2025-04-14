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
	SSLMode  string `mapstructure:"DB_SSLMODE"` // new field for SSL mode
}

type NotificationConfig struct {
	Region        string `mapstructure:"AWS_REGION"`
	AccessKey     string `mapstructure:"SNS_ACCESS_KEY"`
	SecretKey     string `mapstructure:"SNS_SECRET_KEY"`
	SessionToken  string `mapstructure:"SNS_SESSION_TOKEN"`
	AdminTopicArn string `mapstructure:"ADMIN_TOPIC_ARN"`
	UserTopicArn  string `mapstructure:"USER_TOPIC_ARN"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

// Config is the main configuration struct for your application.
type Config struct {
	DB            DBConfig           `mapstructure:",squash"`
	RedisConfig   RedisConfig        `mapstructure:",squash"`
	Notification  NotificationConfig `mapstructure:",squash"`
	AppEnv        string             `mapstructure:"APP_ENV"`
	ServerAddress string             `mapstructure:"SERVER_ADDRESS"`
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
