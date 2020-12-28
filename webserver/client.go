package webserver

import (
	"encoding/json"
	"log"
	"math/rand"

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

func (c *Client) ReadMsg() {
	// Reads msg from the user and sends it to the hub
	// Does the security checks
	// Will edit the message to add the next color
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	h := c.Hub
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		newInfo := new(WSData)
		if err := json.Unmarshal(msg, newInfo); err != nil {
			log.Println(err)
			return
		}
		switch newInfo.Type {
		case "start":
			// Person wants to start the game
			if c.Leader {
				/* We will only see "start" in beginning of each game
				Make Color list / order of player moves
				iterating through map is already random
				*/
				newInfo.Rows = makeLegal(newInfo.Rows)
				newInfo.Cols = makeLegal(newInfo.Cols)
				h.Colors = make([]string, len(h.Clients))
				index := 0
				for client, _ := range h.Clients {
					h.Colors[index] = client.Color
					index += 1
				}
				// Shuffle for good measure
				rand.Shuffle(len(h.Colors), func(i, j int) {
					h.Colors[i], h.Colors[j] = h.Colors[j], h.Colors[i]
				})
				next := h.Colors[h.i]
				newInfo.Next = next
				newInfo.Color = next
				newMsg, err := json.Marshal(newInfo)
				if err != nil {
					// Problems in the code
					log.Fatal(err)
					return
				}
				h.Broadcast <- newMsg
			}

		case "move":
			// Handle User move
			if h.Match.IsLegal(newInfo.X, newInfo.Y) {
				ani, static := h.Match.MovePiece(newInfo.X, newInfo.X, c.Color)
				newInfo.Animation = ani
				newInfo.Static = static
				next := h.Colors[h.i]
				newInfo.Next = next
				newInfo.Color = next
				h.i += 1
				if h.i == len(h.Colors) {
					h.i = 0
				}
				newMsg, err := json.Marshal(newInfo)
				if err != nil {
					// Problems in the code
					log.Fatal(err)
					return
				}
				h.Broadcast <- newMsg
			}
		}
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
