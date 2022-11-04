package config

import (
	"fmt"
	"log"
	"os"

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

// https://www.velotio.com/engineering-blog/build-a-containerized-microservice-in-golang

var (
	DB *gorm.DB
	err error
)

func Setup() *gorm.DB {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Host,
		Port,
		User,
		Name,
		Password,
	)

	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
	  SkipDefaultTransaction: true, // for performance reasons
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
        os.Exit(2)
	}
	log.Println("Database connected")

	// when you want to migrate to db uncomment this
	return DB
}
