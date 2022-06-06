package main

import (
	"fmt"
	"log"
	"errors"
	"gorm.io/gorm"

	"example.com/chatapp/db/config"
	"example.com/chatapp/db/models"
	// "example.com/chatapp/db/utils"
)

func main() {
	// connect to postgresql

	db, err := config.Setup()

	if err != nil {
		log.Panic(err)
		return
	}

	fmt.Println("Connected.")

	// migrate models
	db.AutoMigrate(&models.User{})
	fmt.Println("Successfully migrated.")

	// create user
	// hashedPassword,err := utils.HashPassword("123456789")
	// if err != nil{
	// 	log.Panic("Password does not not hashed")
	// 	return
	// }
	// user := models.User{
	// 	Username: "test2",
	// 	Password: hashedPassword,
	// }
	// result,err := createUser(db,user)
	// if err != nil{
	// 	log.Panic(err)
	// 	return
	// }
	// fmt.Println("User created",result)
}

func createUser(db *gorm.DB, user models.User) (int64, error) {
	result := db.Create(&user)
	fmt.Println(result,result.RowsAffected)
	if result.RowsAffected == 0 {
		return 0, errors.New("User not created")
	}
	return result.RowsAffected, nil
}