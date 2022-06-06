package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

	"example.com/chatapp/utils"
	"example.com/chatapp/db/models"
	dbutils "example.com/chatapp/db/utils"
)

// https://golang.org/doc/code.html

func CreateToken(user models.User) (string, error) {
    var err error

    // Creating Access Token
    atClaims := jwt.MapClaims{}
    atClaims["authorized"] = true
    atClaims["user"] = map[string]any{"id":user.ID,"username": user.Username}
    atClaims["exp"] = time.Now().Add(time.Minute * 90).Unix()

    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    token, err := at.SignedString([]byte(utils.GetDotEnvVariable("ACCESS_SECRET")))
    
    if err != nil {
	    return "", err
    }
    return token, nil
}

func Login(c *gin.Context) {
	var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusUnprocessableEntity, "Invalid json provided.")
		return
	}
	if u.Username == "" || u.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Username or Password not filled.",
		})
		return
	}

	db,err := dbutils.SetupDB()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,gin.H{
			"error": "Something went wrong",
		})
		return
	}

	user,is_registered := dbutils.GetUser(db,u.Username)

	if !is_registered {
		c.IndentedJSON(http.StatusNotFound,gin.H{
			"error": "User not found",
		})
		return
	}

	valid := dbutils.CheckPasswordHash(u.Password,user.Password)

	if !valid {
		c.IndentedJSON(http.StatusNotFound,gin.H{
			"error": "Password is invalid.",
		})
		return
	} 

	token,err := CreateToken(user)

	c.IndentedJSON(http.StatusOK,gin.H{"token": token})

}
