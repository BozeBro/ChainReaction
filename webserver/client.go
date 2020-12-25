package webserver

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Color  string

	Leader bool

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	received chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) readMsg() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c.hub.broadcast <- msg
	}
}
func (c *Client) writeMsg() {
	txtMsg := 1
	for {
		select {
		case msg := <-c.received:
			err := c.conn.WriteMessage(txtMsg, msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
