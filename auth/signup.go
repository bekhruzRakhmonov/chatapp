package auth

import (
	_ "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"example.com/chatapp/db/config"
	"example.com/chatapp/db/models"
	"example.com/chatapp/db/utils"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []User{}

func createUser(db *gorm.DB, user models.User) (int64, error) {
	result := db.Create(&user)
	log.Println(result, *result)
	if result.RowsAffected == 0 {
		return 0, errors.New("User not created")
	}
	return result.RowsAffected, nil
}

func CreateUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Fatal(err)
		return
	}

	db, err := config.Setup()

	if err != nil {
		log.Panic(err)
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		log.Panic("Password does not not hashed")
		return
	}
	user := models.User{
		Username: newUser.Username,
		Password: hashedPassword,
	}

	u,is_registered := utils.GetUser(db, newUser.Username)
	_ = u
	fmt.Println("Error:", is_registered)
	if is_registered == true {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "User has already registered.",
		})
		return
	}

	result, err := createUser(db, user)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("User created", result)

	c.IndentedJSON(http.StatusCreated, newUser)
}
