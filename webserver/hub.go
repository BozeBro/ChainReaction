package webserver

import (
	"encoding/json"
	"log"
)

type Hub struct {
	// Tells server if the Hub is running
	Alive bool
	// Channel that tells server that a player left
	Delete chan bool
	// Channel telling server to remove id / kill the hub
	Stop chan bool
	// Mapping of clients. Unordered
	Clients map[*Client]bool
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
type RoomData struct {
	/*
		Similar to server's GameData but without the chans
	*/
	Room, Pin    string
	Players, Max int
}

func NewHub(roomData *RoomData) *Hub {
	return &Hub{
		Alive:      false,
		Delete:     make(chan bool),
		Stop:       make(chan bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		RoomData:   roomData,
		Match:      new(Chain),
	}
}
func (h *Hub) GetUniqueColor(c string) string {
	for client, _ := range h.Clients {
		if c == client.Color {
			return h.GetUniqueColor(RandomColor())
		}
	}
	return c
}
func (h *Hub) Run() {
	/*
		Equivalent to turning on the computer
		Handles the registering, unregistering, and broadcasting
		Will kill itself when all the players leave

	*/
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
			colorJson := &struct {
				Color string
				Type  string
			}{Color: client.Color, Type: "color"}
			payload, err := json.Marshal(colorJson)
			if err != nil {
				// This should never happen.
				// Only in bugs
				log.Fatal(err)
				return
			}
			h.Clients[client] = true
			// Send player info on his color
			client.Received <- payload
			h.RoomData.Players += 1
			h.Update()
		case client := <-h.Unregister:
			delete(h.Clients, client)
			close(client.Received)
			h.Delete <- true
			if len(h.Clients) == 0 {
				return
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
func (h *Hub) Update() {
	/*
		Function tells front end how many players are in the lobby
	*/
	players := struct {
		Type    string `json:"type"`
		Players int    `json:"players"`
	}{
		Type:    "update",
		Players: h.RoomData.Players,
	}
	payload, err := json.Marshal(players)
	if err != nil {
		// Error should only happen if a bug is here
		log.Fatal(err)
		return
	}
	for client := range h.Clients {
		select {
		case client.Received <- payload:
		default:
			delete(h.Clients, client)
			close(client.Received)
		}
	}
}
