package websocket

import (
	"encoding/json"
	"log"
	"math/rand"
)

// Type that will deal with msgs from websocket
type Responder func(*WSData) error

func (c *Client) start(game Game) Responder {
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
		// Current person's turn
		// used 0 here to hammer in that 0 is the first person
		// reset h.i fior when game is restarted
		h.i = 0
		playInfo.Turn = h.Colors[h.i]
		h.Match.InitBoard(playInfo.Rows, playInfo.Cols)
		newMsg, err := json.Marshal(playInfo)
		if err != nil {
			// Problems in the code
			// Try again
			return err
		}
		h.Broadcast <- newMsg
		return nil
	})
}

func (c *Client) move() Responder {
	h := c.Hub
	return Responder(func(playInfo *WSData) error {
		// See if a person can click the square or not.
		// Within bounds and compatible color
		isLegal := h.Match.IsLegalMove(playInfo.X, playInfo.Y, c.Color)
		if isLegal && c.Color == h.Colors[h.i] {
			// Move Piece, Update colorMap, record animation and new positions
			ani, static := h.Match.MovePiece(playInfo.X, playInfo.Y, c.Color)
			// Color Slice will be updated in the MovePiece Functions
			h.i++
			if h.i >= len(h.Colors) {
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
				return err

			}
			h.Broadcast <- newMsg
		}
		if len(h.Colors) == 1 {
			for _, v := range h.Clients {
				log.Print(v)
			}
			// The last player is declared the winner
			err := h.end(h.Colors[0])
			if err != nil {
				return err
			}
			h.Colors = make([]string, 0, 5)
		}
		return nil
	})

}
