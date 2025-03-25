package db

import (
	"fmt"
	"time"

	"github.com/guncv/Poll-Voting-Website/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB attempts to initialize the database connection using the provided configuration.
// It will try up to 5 times before returning nil if the connection cannot be established.
func InitDB(cfg config.Config) *gorm.DB {
	var err error

	dbSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port,
	)

	for i := 1; i <= 5; i++ {
		DB, err = gorm.Open(postgres.Open(dbSource), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to database successfully")
			return DB
		}

		fmt.Printf("Attempt %d: Failed to connect to database: %v. Retrying in 2 seconds...\n", i, err)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Failed to connect to database after 5 attempts")
	return nil
}
