package chat

import (
	"log"
	"errors"
	_"fmt"
	_"time"
	"gorm.io/gorm"

	"example.com/chatapp/db/models"
	dbutils "example.com/chatapp/db/utils"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	maxUsersCount uint8
}

func NewHub(maxUsersCount uint8) *Hub {
	return &Hub{
		broadcast:  		make(chan *Message),
		register:   		make(chan *Client),
		unregister: 		make(chan *Client),
		clients:    		make(map[*Client]bool),
		maxUsersCount: 		maxUsersCount,
	}
}

func setupDB() (*gorm.DB,error) {
	db,err := dbutils.SetupDB()
	if err != nil {
		return db, errors.New("Database is not connected")
	}
	return db,nil
} 

func IsUserRegistered(outbound_user map[string]any,client *Client) bool {
	db,err := setupDB()
	if err != nil {
		client.send <- &Message{Outbound: client.peer.outbound, Inbound: client.peer.inbound, Send: &ChildMesssage{Status: 500, Message: "Internal server error."}}// []byte("{\"status\": 500,\"error\":\"Internal server error.\"}")
		close(client.send)
		return false
	}
	from,_ := dbutils.GetUser(db,client.peer.outbound)
	outbound_username,ok := outbound_user["username"].(string)

	if ok {
		to,_ := dbutils.GetUser(db,outbound_username)
		log.Println(from,to)
		return true
	}
	return false

}

func RegisterUser(outbound_user map[string]any,client *Client) {
	db,err := setupDB()
	if err != nil {
		client.send <- &Message{Outbound: client.peer.outbound, Inbound: client.peer.inbound, Send: &ChildMesssage{Status: 500, Message: "Internal server error."}} // []byte("{\"status\": 500,\"error\":\"Internal server error.\"}")
		close(client.send)
		return
	}

	from,_ := dbutils.GetUser(db,client.peer.outbound)
	outbound_username,ok := outbound_user["username"].(string)

	if ok {
		to,_ := dbutils.GetUser(db,outbound_username)
		chat := &models.Message{
			From: from,
			To: to,
		}

		rows_affected, err := CreateChat(db,chat)
		_ = rows_affected

		if err != nil {
			return
		}

		log.Println("[hub.go] User registered successfully.")
	} else {
		client.send <- &Message{Outbound: client.peer.outbound, Inbound: client.peer.inbound, Send: &ChildMesssage{Status: 400, Message: "Username is not a string"}}// []byte("{\"status\": 400,\"error\":\"Username is not string.\"}")
	}

}

func CreateChat(db *gorm.DB, chat *models.Message) (int64, error) {
	result := db.Create(chat)
	if result.RowsAffected == 0 {
		return 0, errors.New("Chat is not created")
	}
	return result.RowsAffected, nil
}

func (h *Hub) Run() {
	log.Println("[hub.go] Run() function is started working")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println("Message received:",message.Outbound,message.Inbound,message.Send.Status,message.Send.Message)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}