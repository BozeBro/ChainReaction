package websocket

import (
	"bytes"
	"encoding/json"
	"math/rand"
)

// function type that will deal with a specific msg from websocket
// function should be called by the type they handle
type Responder func(*WSData) error

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// returns a function that handles json data of type start.
// This function should only be used in the beginning of each game.
// Randomizes the player order.
// Sends response of person's turn
// game parameter defines that type of game being played
func (c *Client) start(game *Chain) Responder {
	h := c.Hub
	return Responder(func(playInfo *WSData) error {
		/* We will only see "start" in beginning of each game
		Make Color list / order of player moves
		iterating through map is already random
		*/
		if !c.Leader {
			return nil
		}
		h.Match = game
		h.Colors = make([]string, len(h.Clients))
		index := 0
		for client := range h.Clients {
			h.Colors[index] = client.Color
			h.Clients[client] = 0
			index++
		}
		rand.Shuffle(len(h.Colors), func(i, j int) {
			h.Colors[i], h.Colors[j] = h.Colors[j], h.Colors[i]
		})
		// reset h.i for when game is restarted.
		h.i = 0
		playInfo.Turn = h.Colors[h.i]
		h.Match.InitBoard(playInfo.Rows, playInfo.Cols)
		for client := range h.Clients {
			if client.Color == h.Colors[0] {
				playInfo.Username = client.Username
			}
		}
		payload, err := json.Marshal(playInfo)
		if err != nil {
			// Problems in the code
			// Try again
			return err
		}
		h.Broadcast <- payload
		return nil
	})
}

// Function handles when a person moves
// Utilizes the Game interface to handle game logic.
// Sends response of animation data and new turn
func (c *Client) Move() Responder {
	h := c.Hub
	return Responder(func(playInfo *WSData) error {
		// See if a person can click the square or not.
		// Within bounds and compatible color
		isLegal := h.Match.IsLegalMove(playInfo.X, playInfo.Y, c.Color)
		if isLegal && c.Color == h.Colors[h.i] {
			// Move Piece, Update colorMap, record animation and new positions
			ani, static := h.Match.MovePiece(playInfo.X, playInfo.Y, c.Color)
			h.i++
			if h.i >= len(h.Colors) {
				h.i = 0
			}
			playInfo.Animation = ani
			playInfo.Static = static
			playInfo.Turn = h.Colors[h.i]
			for client := range h.Clients {
				if client.Color == h.Colors[h.i] {
					playInfo.Username = client.Username
				}
			}
			payload, err := json.Marshal(playInfo)
			if err != nil {
				// Problems in the code
				return err

			}
			h.Broadcast <- payload
		}
		// We have a winner!
		if len(h.Colors) == 1 {
			// The last player is declared the winner
			err := h.End(h.Colors[0])
			if err != nil {
				return err
			}
			// reset colors
			// Cap of 25 because that is max players allowed in a game.
			h.Colors = make([]string, 0, 25)
		}
		return nil
	})

}

func (c *Client) chat() Responder {
	h := c.Hub
	return Responder(func(playInfo *WSData) error {
		msg := string(bytes.TrimSpace(bytes.Replace([]byte(playInfo.Message), newline, space, -1)))
		payload, err := json.Marshal(&struct {
			Type     string `json:"type"`
			Message  string `json:"message"`
			Color    string `json:"color"`
			Username string `json:"username"`
		}{
			Type:     "chat",
			Message:  msg,
			Color:    c.Color,
			Username: c.Username,
		})
		if err != nil {
			return err
		}
		h.Broadcast <- payload
		return nil
	})
}
