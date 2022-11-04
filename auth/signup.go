package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	db "example.com/chatapp/db/config"
	"example.com/chatapp/db/models"
	dbutils "example.com/chatapp/db/utils"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) CheckNotBlank(c *gin.Context) {
	if u.Username == "" || u.Password == "" {
		c.IndentedJSON(http.StatusBadRequest,gin.H{
			"error": "Username or Password not filled",
		})
		return
	}
}

func (u *User) Validate(uname,psw string) (string,string) {
	log.Println("User:",u.Username,u.Password)
	username := strings.TrimSpace(uname)
	password := strings.TrimSpace(psw)

	return username, password
}

func createUser(user models.User) (int64, error) {
	result := db.DB.Create(&user)
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

	newUser.CheckNotBlank(c)

	username,password := newUser.Validate(newUser.Username,newUser.Password)

	_,is_registered := dbutils.GetUser(username)
	if is_registered == true {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "User has already registered.",
		})
		return
	}

	start := time.Now()
	hashedPassword, err := dbutils.HashPassword(password)
	if err != nil {
		log.Panic("Password does not not hashed")
		return
	}
	end := time.Since(start)
	log.Println("Time:",end)
	
	user := models.User{
		Username: newUser.Username,
		Password: hashedPassword,
	}

	result, err := createUser(user)
	_ = result
	if err != nil {
		log.Println(err)
		return
	}
	accessToken,refreshToken := GenerateTokens(user.Username)

	c.IndentedJSON(http.StatusCreated,gin.H{"access_token": accessToken,"refresh_token":refreshToken})
}
