package webserver

import (
	"encoding/json"
	"log"
)

// Hub is the server representative that provides a link to the clients
// Handles stopping itself, tracking the players, keeping data, broadcasting, registering,
// unregistering
type Hub struct {
	// Tells server if the Hub is running
	Alive bool
	// Channel telling server to remove id / kill the hub
	Stop chan bool
	// Mapping of clients. Unordered
	Clients map[*Client]int
	// Channel to tell the http that a player left
	Leaver chan bool
	// incoming broadcasting reqs from clients
	Broadcast chan []byte
	// Register requests from the clients.
	Register chan *Client
	// Unregister requests from clients.
	Unregister chan *Client
	// Index of who's turn it is
	i int
	// Move Order of players
	Colors []string
	// Has the data of the room
	RoomData *RoomData
	// Tracker of Game State. "Match" name to not confuse namespace
	Match Game
}

// RoomData provides about the. Will be for players trying to enter the room
type RoomData struct {
	/*
		Similar to server's GameData but without the chans
	*/
	Room, Pin    string
	Players, Max int
}

// NewHub Creates a newHub for a game to take place in
func NewHub(roomData *RoomData) *Hub {
	return &Hub{
		Alive:      false,
		Stop:       make(chan bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]int),
		RoomData:   roomData,
	}
}

// GetUniqueColor grabs a unique from COLORS in utility.go
// It makes sure the color is unique
func (h *Hub) GetUniqueColor(c string) string {
	for client := range h.Clients {
		if c == client.Color {
			return h.GetUniqueColor(RandomColor())
		}
	}
	return c
}

//  Run is equivalent to turning on the computer
//	Handles the registering, unregistering, and broadcasting
//	Will kill itself when all the players leave
func (h *Hub) Run() {
	defer func() {
		h.Stop <- true
	}()
	// No one can join the game
	h.Alive = true
	// wait for register, unregister, or broadcast chan to be filled
	for {
		select {
		case client := <-h.Register:
			// Assign unique color
			client.Color = h.GetUniqueColor(RandomColor())
			colorJSON := &struct {
				Color string `json:"color"`
				Type  string `json:"type"`
			}{Color: client.Color, Type: "color"}
			payload, err := json.Marshal(colorJSON)
			if err != nil {
				// This should never happen.
				// Only in bugs
				log.Println(err)
				return
			}
			h.Clients[client] = 0
			// Tell user what color the person is
			go func() { client.Received <- payload }()
			h.RoomData.Players++
			// Update the amount of players in the lobby
			go h.Update()
		case client := <-h.Unregister:
			delete(h.Clients, client)
			close(client.Received)
			if len(h.Clients) == 0 {
				// NO one is in the lobby
				return
			}
			h.RoomData.Players--
			go h.Update()
			go func() {
				h.Leaver <- true
			}()
			if len(h.Clients) == 1 && len(h.Colors) > 0 {
				// The alone player is the winner
				for client := range h.Clients {
					// must loop to get the person
					err := h.end(client.Color)
					if err != nil {
						log.Println(err)
						return
					}
				}

			} else if len(h.Colors) > 0 {
				// Handle if leaver was its turn
				curTurn := h.Colors[h.i]
				for index, color := range h.Colors {
					if client.Color == color {
						h.Colors = append(h.Colors[:index], h.Colors[index+1:]...)
					}
				}
				if h.i > len(h.Colors) {
					h.i = 0
				}
				if curTurn != h.Colors[h.i] {
					payload := &WSData{
						Turn: h.Colors[h.i],
						Type: "changeColor",
					}
					newMsg, err := json.Marshal(payload)
					if err != nil {
						log.Println(err)
						return
					}
					go func() {
						h.Broadcast <- newMsg
					}()
				}
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Received <- message:
				default:
					delete(h.Clients, client)
					close(client.Received)
				}
			}
		}
	}
}

// Update tells front end how many players are in the lobby
func (h *Hub) Update() {
	players := &struct {
		Type    string `json:"type"`
		Players int    `json:"players"`
	}{
		Type:    "update",
		Players: h.RoomData.Players,
	}
	payload, err := json.Marshal(players)
	if err != nil {
		// Error should only happen if a bug is here
		log.Println(err)
		return
	}
	h.Broadcast <- payload
}

func (h *Hub) end(color string) error {
	payload := &struct {
		Type   string `json:"type"`
		Winner string `json:"winner"`
	}{Type: "end", Winner: color}
	msg, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	go func() {
		h.Broadcast <- msg
	}()
	return nil
}
func (h *Hub) CloseChans() {
	close(h.Stop)
	close(h.Broadcast)
	close(h.Register)
	close(h.Unregister)
}
