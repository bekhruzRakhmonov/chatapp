package main

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"

	"example.com/chatapp/auth"
	"example.com/chatapp/middleware"
	"example.com/chatapp/chat"
	"example.com/chatapp/searchapp"
)

func main() {
	r := gin.Default()

	r.POST("/signup", middleware.CheckRequestHeader(), auth.CreateUser)
	r.GET("/signup",func (c *gin.Context) {
		c.IndentedJSON(http.StatusMethodNotAllowed,gin.H{
			"error": "Method not allowed.",
		})
	})
	r.POST("/login", middleware.CheckRequestHeader(), auth.Login)
	r.GET("/login",func (c *gin.Context) {
		c.IndentedJSON(http.StatusMethodNotAllowed,gin.H{
			"error": "Method not allowed.",
		})
	})

	r.GET("/test",middleware.EnsureLoggedIn(), func (c *gin.Context){
		// props, ok := c.Request.Context().Value("props").(jwt.MapClaims)

		if props,exists := c.Get("props"); exists {
			log.Println(props)
			c.IndentedJSON(http.StatusOK,"User is logged in")
		} else {
			c.IndentedJSON(http.StatusUnauthorized,gin.H{"error": "Unauthorized user."})
		}

	})

	hub := chat.NewHub(2)
	go hub.Run()
	r.GET("/ws/:username", middleware.EnsureLoggedIn(), func(c *gin.Context) {
		authorized,_ := c.Get("is_authorized")
		is_authorized,_ := authorized.(bool)
		// c.Request.Host+c.Request.URL.Path
		username := c.Param("username")

		chat.ServeWs(hub, c, username,is_authorized)
	})

	r.GET("/search", middleware.EnsureLoggedIn(), func(c *gin.Context) {
		authorized,ok := c.Get("is_authorized")

		is_authorized := false
		if ok {
			is_authorized,_ = authorized.(bool)
		}
		searchapp.Run(c,is_authorized)
	})

	// https://github.com/egaprsty/ChatAppGolang
	
	r.Run(":8000")
}
