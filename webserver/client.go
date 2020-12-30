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

// WSData provides allowed fields to be received from the front end
type WSData struct {
	Type      string    `json:"type"` // Type of message allows front end to know how to deal with it
	X         int       `json:"x"`    // X coordinate clicked
	Y         int       `json:"y"`    // Y coordinate clicked
	Turn      string    `json:"turn"` // players turn
	Rows      int       `json:"rows"`
	Cols      int       `json:"cols"`
	Animation [][][]int `json:"animation"` // Instrutios on animation
	Static    [][][]int `json:"static"`    // What the new board will look like
}

// ReadMsg Reads msg from the user and sends it to the hub
// Does the security checks
// Will edit the message to add the next color
func (c *Client) ReadMsg() {
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
		playInfo := new(WSData)
		if err := json.Unmarshal(msg, playInfo); err != nil {
			log.Println(err)
			return
		}
		switch playInfo.Type {
		case "start":
			// Person wants to start the game
			if c.Leader {
				/* We will only see "start" in beginning of each game
				Make Color list / order of player moves
				iterating through map is already random
				*/
				h.Match = &Chain{Hub: h}
				h.Colors = make([]string, len(h.Clients))
				index := 0
				for client := range h.Clients {
					h.Colors[index] = client.Color
					index++
				}
				// Randomize players
				rand.Shuffle(len(h.Colors), func(i, j int) {
					h.Colors[i], h.Colors[j] = h.Colors[j], h.Colors[i]
				})
				// Current person's turn
				playInfo.Turn = h.Colors[0]
				h.Match.InitBoard(playInfo.Rows, playInfo.Cols)
				newMsg, err := json.Marshal(playInfo)
				if err != nil {
					// Problems in the code
					log.Fatal(err)
					break
				}
				h.Broadcast <- newMsg
			}
		case "move":
			// See if a person can click the square or not.
			// Within bounds and compatible color
			isLegal := h.Match.IsLegalMove(playInfo.X, playInfo.Y, c.Color)
			if isLegal {
				// Move Piece, Update colorMap, record animation and new positions
				ani, static := h.Match.MovePiece(playInfo.X, playInfo.Y, c.Color)
				// Color Slice will be updated in the MovePiece Functions
				if len(h.Colors) > 1 {
					h.i++
					if h.i == len(h.Colors) {
						h.i = 0
					}
					// Game is not over yet
					playInfo.Animation = ani
					playInfo.Static = static
					// Getting next person
					playInfo.Turn = h.Colors[h.i]
					newMsg, err := json.Marshal(playInfo)
					if err != nil {
						// Problems in the code
						log.Fatal(err)
						break
					}
					h.Broadcast <- newMsg

				} else {
					// One player remaining. That player wins.
					log.Println("WINNER WINNER")
				}
			}
		}
	}
}

// WriteMsg sends msg from the hub to the client
func (c *Client) WriteMsg() {
	txtMsg := 1
	for {
		select {
		case msg := <-c.Received:
			err := c.Conn.WriteMessage(txtMsg, msg)
			if err != nil {
				return
			}
		}
	}
}
