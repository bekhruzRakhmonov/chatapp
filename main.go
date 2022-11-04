package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/chatapp/auth"
	"example.com/chatapp/chat"
	db "example.com/chatapp/db/config"
	_ "example.com/chatapp/db/models"
	"example.com/chatapp/middleware"
	"example.com/chatapp/searchapp"
)

func main() {
	db.Setup()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	account := r.Group("/account")
	{
		account.POST("/signup", auth.CreateUser)
		account.GET("/signup", func(c *gin.Context) {
			c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Method not allowed.",
			})
		})
		account.POST("/login", middleware.CheckRequestHeader(), auth.Login)
		account.GET("/login", func(c *gin.Context) {
			c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Method not allowed.",
			})
		})

		account.POST("/get-access-token", auth.GetAccessToken)
	}
	r.GET("/test", middleware.EnsureLoggedIn(), func(c *gin.Context) {
		//props, ok := c.Request.Context().Value("props").(*models.Claims{})

		if props, exists := c.Get("props"); exists {
			log.Println(props)
			c.IndentedJSON(http.StatusOK, "User is logged in")
		} else {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user."})
		}

	})

	chatGroup := r.Group("/chat")
	{
		hub := chat.NewHub()
		go hub.Run()
		chatGroup.GET(":username", func(c *gin.Context) {
			is_authorized := chat.EnsureLoggedIn(c)
			log.Println("user is authorized", is_authorized)
			username := c.Param("username")

			log.Println("Username is ", username)

			chat.ServeWs(hub, c, username, is_authorized)

			// c.Request.Host+c.Request.URL.Path
		})
		chatGroup.GET("/get-chats", middleware.EnsureLoggedIn(), chat.GetChats)
		chatGroup.GET("/get-last-chat", middleware.EnsureLoggedIn(), chat.GetLastChat)

		chatGroup.GET("/get-messages/:inbound", middleware.EnsureLoggedIn(), chat.GetMessages)

		chatGroup.GET("/search", func(c *gin.Context) {
			is_authorized := chat.EnsureLoggedIn(c)
			searchapp.Run(c, is_authorized)
		})
	}

	r.Run(":8000")
}

// https://github.com/egaprsty/ChatAppGolang
