package middleware

import (
	"net/http"
	"log"
	"fmt"
	"strings"
	"time"
	_ "context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"

	"example.com/chatapp/utils"
	"example.com/chatapp/db/models"
)

// https://chenyitian.gitbooks.io/gin-tutorials/content/tdd/20.html

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Headers","Content-Type,Authorization")

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
			claims := new(models.Claims)
			token, err := jwt.ParseWithClaims(jwtToken, claims, func (token *jwt.Token) (interface{},error){
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New(fmt.Sprintf("Unexpected signing method %s",token.Header["alg"]))
				}
				return []byte(utils.GetDotEnvVariable("ACCESS_SECRET")),nil
			})

			if token.Valid {
				log.Println(claims.ExpiresAt < time.Now().Unix(),claims.ExpiresAt,time.Now().Unix())
				if claims.ExpiresAt < time.Now().Unix() {
	                c.IndentedJSON(http.StatusBadRequest,gin.H{"error":"Token has expired."})
	            } else {
	            	c.Set("is_authorized",true)
					c.Set("props",claims)
	            }
	            return
			} else if ve, ok := err.(*jwt.ValidationError); ok {
	            if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	                // this is not even a token, we should delete the cookies here
	                // c.ClearCookie("access_token", "refresh_token")
	                // c.IndentedJSON(http.StatusForbidden,gin.H{"error":"Unauthorized user."})
	                c.Set("is_authorized",false)
	                return
	            } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	                // Token is either expired or not active yet
	                c.Set("is_authorized",false)
	                return
	            }
            }

			// if claims, ok := token.Claims.(claims{}); ok && token.Valid {
			// 	// I can use ctx with builtin http/net server 
 		// 		// ctx := context.WithValue(c.Request.Context(),"props",claims)
 		// 		// Access context values in handlers like this
 		// 		// props, _ := r.Context().Value("props").(jwt.MapClaims)
 		// 		// log.Println("Context:",ctx)
			// 	c.Set("is_authorized",true)
			// 	c.Set("props",claims)
			// } else {
			// 	log.Println(err)
			// 	c.Set("is_authorized",false)
			// }
		}
		c.Next()
	}
}


func EnsureNotLoggedIn() gin.HandlerFunc {

	return func (c *gin.Context) {
		authHeader := strings.Split(c.Request.Header.Get("Authorization"),"Bearer ")
		if len(authHeader) > 1 {
			c.IndentedJSON(http.StatusOK,gin.H{
				"error": "You have already authorized.",
			})
		}
	}
}

func CheckRequestHeader() gin.HandlerFunc {
	return func (c *gin.Context) {
		contentType := c.Request.Header["Content-Type"]
		if contentType[0] != "application/json" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": "Content-Type does not provided.",
			})
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}