package chat

import (
	"net/http"

	db "example.com/chatapp/db/config"
	"example.com/chatapp/db/models"

	"github.com/gin-gonic/gin"
)

// `chat/get-messages/:username`
func GetMessages(c *gin.Context) {
	authorized, _ := c.Get("is_authorized")
	is_authorized, _ := authorized.(bool)

	props, exists := c.Get("props")
	claims, ok := props.(*models.Claims)

	if !is_authorized || !exists || !ok {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized user",
		})
		return
	}

	inbound := c.Param("inbound")
	messages := getMessages(claims.Issuer, inbound)

	c.IndentedJSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func getMessages(outbound, inbound string) []models.Message {
	messages := []models.Message{}
	db.DB.Where("(outbound = ? or outbound = ?) and (inbound = ? or inbound = ?)",
		outbound,
		inbound,
		inbound,
		outbound).Preload("Chats").Find(&messages)
	return messages
}

func CreateMessage(outbound, inbound, msg string) {
	chat := CreateChat(outbound, inbound)
	message := &models.Message{
		Outbound: outbound,
		Inbound:  inbound,
		Message:  msg,
	}

	db.DB.Create(message)

	db.DB.Model(chat).Association("Messages").Append(message)
	getMessages(outbound, inbound)

}
