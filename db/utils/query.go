package utils

import (
	"fmt"

	"example.com/chatapp/db/models"
	db "example.com/chatapp/db/config"
)

func GetUser(username string) (models.User,bool) {
	var user models.User
	result := db.DB.First(&user, "username = ?", username)
	if result.RowsAffected == 0 {
		return user,false
	}
	return user,true
}


func FindUser(username string) []string {
	users := []models.User{}
	var usernames []string
	query := fmt.Sprintf("SELECT * FROM users WHERE username LIKE '%s%%'",username)
	db.DB.Raw(query).Scan(&users)

	for _,user := range users {
		usernames = append(usernames,user.Username)
	}

	return usernames
}