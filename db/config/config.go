package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	Host     = "localhost"
	User     = "postgres"
	Password = "Idonotknow1@"
	Name     = "chatapp"
	Port     = "5555"
)

func Setup() (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Host,
		Port,
		User,
		Name,
		Password,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
