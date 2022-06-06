package utils

import (
	"log"
	"errors"
	"gorm.io/gorm"

	"example.com/chatapp/db/config"
)

func SetupDB() (*gorm.DB,error){
	db, err := config.Setup()

	if err != nil {
		log.Panic(err)
		return db,errors.New("Database connection failed.")
	}
	return db,nil
}
