package webserver

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The color that represents the player and on the board
	Color  string
	// Can the player start the game or not?
	Leader bool
	// The hub in which clients will play
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Received chan []byte
}

func (c *Client) ReadMsg() {
	// Reads msg from the user and sends it to the hub
	defer func() {
		c.Hub.Unregister <- c
		c.Hub.Delete <- true
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c.Hub.Broadcast <- msg
	}
}
func (c *Client) WriteMsg() {
	// Sends msg from the hub to the client
	txtMsg := 1
	for {
		select {
		case msg := <-c.Received:
			err := c.Conn.WriteMessage(txtMsg, msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
