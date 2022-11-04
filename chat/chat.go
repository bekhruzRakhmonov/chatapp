package chat

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	db "example.com/chatapp/db/config"
	"example.com/chatapp/db/models"
)

func getAuthenticatedUser(c *gin.Context) (bool, string) {
	authorized, _ := c.Get("is_authorized")
	is_authorized, _ := authorized.(bool)

	props, exists := c.Get("props")
	claims, ok := props.(*models.Claims)

	if is_authorized && exists && ok {
		return true, claims.Issuer
	}
	return false, claims.Issuer
}

// this function responsible for /get-last-chat endpoint
func GetLastChat(c *gin.Context) {
	is_authenticated, user := getAuthenticatedUser(c)

	if !is_authenticated {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	chat, err := getLastChat(user)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	}

	fmt.Println("Chat:", chat)

	if chat.Outbound == user {
		c.IndentedJSON(http.StatusOK, gin.H{
			"success":  "Found",
			"username": chat.Inbound,
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":  "Found",
		"username": chat.Outbound,
	})
}

func GetChat(outbound, inbound string) bool {
	result := db.DB.Where("outbound = ? AND inbound = ?", outbound, inbound).Find(&models.Chat{})

	return result.RowsAffected == 0
}

func GetChats(c *gin.Context) {
	is_authenticated, username := getAuthenticatedUser(c)

	if !is_authenticated {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	chats := getChats(username)

	c.IndentedJSON(http.StatusOK, gin.H{
		"chats": chats,
	})
}

func getChats(username string) []models.Chat {
	chats := []models.Chat{}
	// query := fmt.Sprintf("SELECT * FROM chats WHERE outbound LIKE '%s%%'", username)
	db.DB.Where("outbound = ? or inbound = ?", username, username).Find(&chats)
	return chats
}

func getLastChat(username string) (models.Chat, error) {
	var chat models.Chat
	result := db.DB.First(&chat, "outbound = ?", username)
	if result.RowsAffected == 0 {
		return chat, errors.New("last chat not found")
	}
	return chat, nil
}

func CreateChat(outbound, inbound string) *models.Chat {
	chat := &models.Chat{
		Outbound: outbound,
		Inbound:  inbound,
	}

	db.DB.FirstOrCreate(chat)
	return chat
}
