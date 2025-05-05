package database

import (
	"log"

	"github.com/omarisadev/go-auth/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() {
	var err error
	if DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{}); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected")
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func CreateUser(user *models.User) error {
	result := DB.Create(user)
	return result.Error
}
