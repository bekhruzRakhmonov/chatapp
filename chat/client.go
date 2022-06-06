package chat

import (
	"bytes"
	_"encoding/gob"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	_"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/dgrijalva/jwt-go"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


// peer chat
type Peer struct {
	outbound string
	inbound  string
}

type ChildMesssage struct {
	Status	uint
	Message string
}

type Message struct {
	Outbound string
	Inbound  string
	Send     *ChildMesssage
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *Message

	peer *Peer
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	log.Println("readPump")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		stringed_message := fmt.Sprintf("%s",message)
		c.hub.broadcast <- &Message{Outbound: c.peer.outbound, Inbound: c.peer.inbound, Send: &ChildMesssage{Status: 200, Message: stringed_message}}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	log.Println("writePump")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			log.Println("<-c.send:",message)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			 // var network bytes.Buffer        // Stand-in for a network connection
			 //    encoder := gob.NewEncoder(&network) // Will write to network.
			 //    // decoder := gob.NewDecoder(&network) // Will read from network.

			 //    if err := encoder.Encode(*message); err != nil {
			 //        log.Fatal("encode error:", err)
			 //    }
			 //    log.Println("Network Bytes:",network.Bytes())
	
			if message.Send.Status == 200 {
				data,_ := json.Marshal(message)
				w.Write(data)
			} else if message.Send.Status == 401 {
				data,_ := json.Marshal(message.Send)
				w.Write(data)
				w.Close()
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				message := <-c.send
				data,_ := json.Marshal(message)
				w.Write(data)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Println("writeWait")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, c *gin.Context,username string,is_authorized bool) {
	w := c.Writer
	r := c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// username := r.URL.Query().Get("username")
	props,_ := c.Get("props")
	data,ok := props.(jwt.MapClaims)

	// outbound user
	user,ok2 := data["user"].(map[string]any)
	outbound_username,_ := user["username"].(string)

	client := &Client{hub: hub, conn: conn, send: make(chan *Message, 256),peer: &Peer{outbound:outbound_username,inbound: username}}
	if !is_authorized {
		log.Println("Unauthorized user")
		client.send <- &Message{Outbound: client.peer.outbound, Inbound: client.peer.inbound, Send: &ChildMesssage{Status: 401, Message: "Unauthorized user"}} // []byte("{\"status\": 401,\"error\":\"Unauthorized user.\"}")
		client.writePump()

		defer client.conn.Close()
	} else {
		if ok {
			if ok2 {
				client.hub.register <- client
				is_registered := IsUserRegistered(user,client)
				if !is_registered {
					RegisterUser(user,client)
				}
			}
		} else {
			client.send <- &Message{Outbound: client.peer.outbound, Inbound: client.peer.inbound, Send: &ChildMesssage{Status: 401, Message: "Unauthorized user"}} // []byte("{\"status\": 401,\"error\":\"Unauthorized user.\"}")
		}
	}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}