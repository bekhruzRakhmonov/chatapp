package chat

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"example.com/chatapp/db/models"
	"example.com/chatapp/utils"
)

func EnsureLoggedIn(c *gin.Context) bool {
	accessToken := c.Request.URL.Query().Get("accessToken")
	decodedToken, err := base64.StdEncoding.DecodeString(accessToken)

	log.Println("=====ERR=====", err, decodedToken)

	if err != nil {
		return false
	}

	accessToken = string(decodedToken)

	log.Println("It is accessToken", accessToken)

	claims := new(models.Claims)
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("Unexpected signing method %s", token.Header["alg"]))
		}
		return []byte(utils.GetDotEnvVariable("ACCESS_SECRET")), nil
	})

	if err == nil {
		if token.Valid {
			if claims.ExpiresAt < time.Now().Unix() {
				return false
			} else {
				log.Println("Authorized user")
				c.Set("props", claims)
				return true
			}
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return false
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return false
			}
		}
	}
	return false
}
