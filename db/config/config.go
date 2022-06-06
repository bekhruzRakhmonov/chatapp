package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/chatapp/utils"
)

var (
	Host     = utils.GetDotEnvVariable("DBHOST")
	User     = utils.GetDotEnvVariable("DBUSER")
	Password = utils.GetDotEnvVariable("DBPASSWORD")
	Name     = utils.GetDotEnvVariable("DBNAME")
	Port     = utils.GetDotEnvVariable("DBPORT")
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
