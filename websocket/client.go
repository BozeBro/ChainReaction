package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The color that represents the player and on the board
	Color string
	// Can the player start the game or not?
	Leader bool
	// The hub in which clients will play
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Received chan []byte
}

// WSData provides allowed fields to be received from the front end
// Some other structs are used as one offs in other places in the code.
type WSData struct {
	Type      string    `json:"type"`      // Type of message allows front end to know how to deal with the data
	X         int       `json:"x"`         // X coordinate clicked - "move"
	Y         int       `json:"y"`         // Y coordinate clicked - "move"
	Turn      string    `json:"turn"`      // players turn - "move"
	Animation [][][]int `json:"animation"` // Instrutios on animation - "move"
	Static    [][][]int `json:"static"`    // What the new board will look like - "move"
	Rows      int       `json:"rows"`      // Amount of rows - Sent at "start"
	Cols      int       `json:"cols"`      // Amount of columns - Sent at "start"
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 6 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// ReadMsg Reads msg from the user and sends it to the hub
// Does the security checks and game checking
func (c *Client) ReadMsg() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	resMap := make(map[string]Responder)
	resMap["start"] = c.start(&Chain{Hub: c.Hub})
	resMap["move"] = c.move()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		playInfo := new(WSData)
		if err := json.Unmarshal(msg, playInfo); err != nil {
			log.Println(err)
			return
		}
		resMap[playInfo.Type](playInfo)
	}
}

// WriteMsg sends msg from the hub to the client
func (c *Client) WriteMsg() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Received:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
