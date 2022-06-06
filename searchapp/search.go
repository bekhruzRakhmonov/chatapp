package searchapp


import (
	"log"
	"time"
	"bytes"
	"fmt"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/dgrijalva/jwt-go"

	//  "example.com/chatapp/chat"
	dbutils "example.com/chatapp/db/utils"
)


const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}



type users struct {
	Usernames []string
}

//var users = []user{}

type result struct {
	Username string
	Users *users
	Status int
	Message string
} 

type Client struct {
	conn *websocket.Conn

	// message struct
	send chan *result

	// get username to send exactly
	username string
}

func (c *Client) reader() {
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func (string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType,message,err := c.conn.ReadMessage()

		log.Println("Message Type is",messageType)		

		if err != nil {
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure){
				log.Printf("error: %v",err)
			}
			break
		}

		log.Println(fmt.Sprintf("%s",message))

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		str := fmt.Sprintf("%s",message)

		// str := bytes.NewBuffer(message).String()

		// search user
		usernames := dbutils.FindUser(str)

		c.send <- &result{Username: c.username, Status:200,Users: &users{Usernames: usernames}} // &chat.Message{Outbound: c.username, Send: &chat.ChildMesssage{Status: 200, Message: str}}
	}
}

func (c *Client) writer() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage,[]byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)

			if err != nil {
				return
			}

			if message.Status == 200 {
				data,_ := json.Marshal(message)
				w.Write(data)
			} else if message.Status == 401 {
				data,_ := json.Marshal(message)
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
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func Run(c *gin.Context,is_authorized bool) {
	w := c.Writer
	r := c.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	props,_ := c.Get("props")

	data,_ := props.(jwt.MapClaims)

	user,_ := data["user"].(map[string]any)
	username,_ := user["username"].(string)

	client := &Client{conn: conn, username: username,send: make(chan *result, 256)}

	if !is_authorized {
		client.send <- &result{Status: 401,Message: "Unauthorized user"}  // &chat.Message{Send: &chat.ChildMesssage{Status: 401, Message: "Unauthorized user"}}
		client.writer()

		defer client.conn.Close()
		return
	}

	go client.reader()
	go client.writer()
}