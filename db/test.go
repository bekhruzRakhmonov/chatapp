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

	db := config.Setup()

	// if err != nil {
	// 	log.Panic(err)
	// 	return
	// }

	fmt.Println("Connected.")

	// migrate models
	db.AutoMigrate(&models.Chat{},&models.Message{},&models.User{},&models.Claims{})
	// log.Println("Successfully migrated.")

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

	// user1,_ := utils.GetUser("test05")
	// user2,_ := utils.GetUser("2")

	// message := &models.Message{
	// 	Message: "Salom 2",
	// }

	// db.Create(message)

	// chat := models.Chat{
	// 	OutboundUsername: "feruz",
	// 	InboundUsername: "bexruz",
	// }

	// var chat models.Chat

	// db.FirstOrCreate(&chat)

	// // err := db.SetupJoinTable(chat, "Messages", message)

	// // if err != nil {
	// // 	log.Println(err)
	// // 	return
	// // }

	// db.Model(chat).Association("Messages").Append(message)
	// res := db.Model(&msg).Association("Chats").Find(&chat)
	// db.Model(&chat).Where("outbound_username = ? and inbound_username = ?", "feruz","bexruz").Association("Languages").Find(&models.Message{})
	// log.Println(chat.Messages,chat)

	log.Println("Successfully migrated")
	// log.Println(chat.Outbound)
}

func createUser(db *gorm.DB, user models.User) (int64, error) {
	result := db.Create(&user)
	fmt.Println(result,result.RowsAffected)
	if result.RowsAffected == 0 {
		return 0, errors.New("User not created")
	}
	return result.RowsAffected, nil
}