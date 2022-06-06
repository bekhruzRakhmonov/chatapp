package utils

import (
	"log"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"example.com/chatapp/db/models"
)

func setupDB() (*gorm.DB,error) {
	db,err := SetupDB()
	if err != nil {
		return db, errors.New("Database is not connected")
	}
	return db,nil
}

func GetUser(db *gorm.DB, username string) (models.User,bool) {
	var user models.User
	result := db.First(&user, "username = ?", username)
	if result.RowsAffected == 0 {
		return user,false
	}
	return user,true
}


func FindUser(username string) []string {
	db,err := setupDB()
	if err != nil {
		log.Println(err)
		return []string{}
	}
	users := []models.User{}
	var usernames []string
	query := fmt.Sprintf("SELECT * FROM users WHERE username LIKE '%s%%'",username)
	db.Raw(query).Scan(&users)

	for _,user := range users {
		usernames = append(usernames,user.Username)
	}

	return usernames
}