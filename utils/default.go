package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gimtwi/go-library-project/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Config struct {
	Database struct {
		Username            string `json:"username"`
		HashedAdminPassword string `json:"hashed_admin_password"`
	} `json:"database"`
}

func loadConfig() (Config, error) {
	var config Config
	configFile, err := os.Open("config.json")

	if err != nil {
		return config, err
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func CreateDefaultAdmin(db *gorm.DB) {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("error loading config:", err)
	}

	var admin types.User

	result := db.Where("username = ?", config.Database.Username).First(&admin)
	if result.Error == gorm.ErrRecordNotFound {

		admin = types.User{
			ID:       uuid.NewString(),
			Username: config.Database.Username,
			Password: config.Database.HashedAdminPassword,
			Role:     types.Admin,
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Fatal("error creating default admin user:", err)
		}
		fmt.Println("default admin user created successfully")
	} else if result.Error != nil {
		log.Fatal("error checking admin user:", result.Error)
	} else {
		fmt.Println("admin user already exists")
	}
}
