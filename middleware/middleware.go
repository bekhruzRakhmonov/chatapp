package middleware

import (
	"net/http"
	"log"
	"fmt"
	"strings"
	_ "context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

	"example.com/chatapp/utils"
)

// https://chenyitian.gitbooks.io/gin-tutorials/content/tdd/20.html

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
	
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// https://hackernoon.com/creating-a-middleware-in-golang-for-jwt-based-authentication-cx3f32z8

func EnsureLoggedIn() gin.HandlerFunc {

	return func (c *gin.Context) {
		authHeader := strings.Split(c.Request.Header.Get("Authorization"),"Bearer ")
		if len(authHeader) != 2 {
			c.Set("is_authorized",false)
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func (token *jwt.Token) (interface{},error){
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New(fmt.Sprintf("Unexpected signing method %s",token.Header["alg"]))
				}
				return []byte(utils.GetDotEnvVariable("ACCESS_SECRET")),nil
			})
		
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// I can use ctx with builtin http/net server 
 				// ctx := context.WithValue(c.Request.Context(),"props",claims)
 				// Access context values in handlers like this
 				// props, _ := r.Context().Value("props").(jwt.MapClaims)
 				// log.Println("Context:",ctx)
				c.Set("is_authorized",true)
				c.Set("props",claims)
			} else {
				log.Println(err)
				c.Set("is_authorized",false)
			}
		}
	}
}

func EnsureNotLoggedIn() gin.HandlerFunc {

	return func (c *gin.Context) {
		authHeader := strings.Split(c.Request.Header.Get("Authorization"),"Bearer ")
		if len(authHeader) != 2 {
			c.IndentedJSON(http.StatusUnauthorized,gin.H{
				"error": "Unauthorized user.",
			})
		}
	}
}

func CheckRequestHeader() gin.HandlerFunc {
	return func (c *gin.Context) {
		contentType := c.Request.Header["Content-Type"]
		log.Println(contentType)
		log.Println(contentType[0])
		if contentType[0] != "application/json" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": "Content-Type does not provided.",
			})
		}
	}
}