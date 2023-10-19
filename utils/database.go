package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/gimtwi/go-library-project/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = gorm.Open(postgres.Open(dbInfo), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	fmt.Println("connected to the database!")
}

func MigrateDB() {
	DB.AutoMigrate(&types.User{}, &types.Book{}, &types.Author{}, &types.Genre{}, &types.Hold{}, &types.Loan{}, &types.Notification{})
	fmt.Println("database migration completed successfully!")
}
