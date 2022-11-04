package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	dbutils "example.com/chatapp/db/utils"
)

// func CreateToken(user models.User) (string, error) {
//     var err error

//     // Creating Access Token
//     atClaims := jwt.MapClaims{}
//     atClaims["authorized"] = true
//     atClaims["user"] = map[string]any{"id":user.ID,"username": user.Username}
//     atClaims["exp"] = time.Now().Add(time.Minute * 90).Unix()

//     at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
//     token, err := at.SignedString([]byte(utils.GetDotEnvVariable("ACCESS_SECRET")))
    
//     if err != nil {
// 	    return "", err
//     }

//     return token, nil
// }

func Login(c *gin.Context) {
	var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{
			"error":"Invalid json provided.",
		})
		return
	}

	u.CheckNotBlank(c)

	username,password := u.Validate(u.Username,u.Password)

	user,is_registered := dbutils.GetUser(username)

	if !is_registered {
		c.IndentedJSON(http.StatusNotFound,gin.H{
			"error": "User not found",
		})
		return
	}

	valid := dbutils.CheckPasswordHash(password,user.Password)

	if !valid {
		c.IndentedJSON(http.StatusBadRequest,gin.H{
			"error": "Password is invalid.",
		})
		return
	} 

	start := time.Now()
	accessToken,refreshToken := GenerateTokens(user.Username)
	end := time.Since(start)
	log.Println("end",end)

	c.IndentedJSON(http.StatusOK,gin.H{"access_token": accessToken,"refresh_token":refreshToken})
}
