package db

import (
	"authentication/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB(connectionString string) error {
	var err error

	// Open a new db connection
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate User model
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	fmt.Println("Connected Successfully to the Database")

	return nil
}
