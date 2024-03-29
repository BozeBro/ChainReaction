package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// Username of the player
	Username string

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

	// channel to kill client. Only Used by botClients
	Stop chan bool
}

// WSData provides allowed fields to be received from the front end
// Some other structs are used as single uses in other places in the code.
type WSData struct {
	Type      string    `json:"type"`      // Type of message allows front end to know how to deal with the data
	X         int       `json:"x"`         // X coordinate clicked - "move"
	Y         int       `json:"y"`         // Y coordinate clicked - "move"
	Turn      string    `json:"turn"`      // players turn - "move"
	Animation [][][]int `json:"animation"` // Instrutios on animation - "move"
	Static    [][][]int `json:"static"`    // What the new board will look like - "move"
	Rows      int       `json:"rows"`      // Amount of rows - Sent at "start"
	Cols      int       `json:"cols"`      // Amount of columns - Sent at "start"
	Message   string    `json:"message"`   // chat messsage sent by a user - "chat"
	Color     string    `json:"color"`     // color of the person, used once - "color"
	Username  string    `json:"username"`  // username of each player
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
	// Mapping that will handle how to respond to a websocket msg
	resMap := make(map[string]Responder)
	resMap["start"] = c.start(&Chain{Hub: c.Hub})
	resMap["move"] = c.Move()
	resMap["chat"] = c.chat()
	for {
		// Close the go routine when hub is closed
		if len(c.Stop) > 0 {
			return
		}
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
    if _, ok := resMap[playInfo.Type]; !ok {
      // Ignore unknown json
      log.Printf("JSON type %s does not exist, ignore", playInfo.Type)
      continue
    }

		if err := resMap[playInfo.Type](playInfo); err != nil {
			log.Println(err)
			return
		}
	}
}

//  WriteMsg sends msg from the hub to the client
//  Contains Ping Handler implementation.
//  See RFC5.5.2 https://tools.ietf.org/html/rfc6455#section-5.5.2 for more info
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
			// Sending Ping msg
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
