package db

import (
	"fmt"
	"time"

	"github.com/guncv/Poll-Voting-Website/backend/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(config util.Config) *gorm.DB {
	var err error

	dbSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort,
	)

	for i := 1; i <= 5; i++ {
		DB, err = gorm.Open(postgres.Open(dbSource), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to database successfully")
			return nil
		}

		fmt.Printf("Attempt %d: Failed to connect to database, retrying in 2s...\n", i)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Database connected")

	return DB
}
