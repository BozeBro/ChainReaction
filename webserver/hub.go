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
	GameData *GameData
}
type GameData struct {
	/*
		Similar to server's GameData but without the chans
	*/
	Room, Pin    string
	Players, Max int
}
type WSData struct {
	Type  string `json:"type"` // Type of message allows front end to know how to deal with it
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
	Val   bool   `json:"val"`
	Next  string `json:"next"` // Next Color
}

func NewHub(gameData *GameData) *Hub {
	return &Hub{
		Alive:      false,
		Delete:     make(chan bool),
		Stop:       make(chan bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		GameData:   gameData,
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
		h.Alive = false
		h.Stop <- true
		return
	}()
	h.Alive = true
	// wait for register, unregister, or broadcast chan to be filled
	for {
		select {
		case client := <-h.Register:
			// Assign unique color
			client.Color = h.GetUniqueColor(RandomColor())
			wsData := &WSData{Color: client.Color, Type: "color"}
			log.Println(client.Color)
			payload, err := json.Marshal(wsData)
			if err != nil {
				// This should never happen.
				// Only in bugs
				log.Fatal(err)
				return
			}
			h.Clients[client] = true
			// Send player info on his color
			client.Received <- payload
			h.GameData.Players += 1
			go h.Update()
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Received)
				h.Delete <- true
				if len(h.Clients) == 0 {
					return
				}
			}
		case message := <-h.Broadcast:
			// EditMsg will add the next color onto the request
			// It will readjust whose next turn it will be
			newMsg := h.EditMsg(message)
			for client := range h.Clients {
				select {
				case client.Received <- newMsg:
				default:
					close(client.Received)
					delete(h.Clients, client)
				}
			}
		}
	}
}
func (h *Hub) EditMsg(msg []byte) []byte {
	// Add the next color to the broadcasted msg
	newInfo := &WSData{}
	if err := json.Unmarshal(msg, newInfo); err != nil {
		log.Fatal(err)
		return nil
	}
	if newInfo.Type == "start" {
		/* We will only see "start" in beginning of each game
		Make Color list / order of player moves
		iterating through map is already random
		*/
		h.Colors = make([]string, len(h.Clients))
		index := 0
		for client, _ := range h.Clients {
			h.Colors[index] = client.Color
			index += 1
		}
	}
	next := h.Colors[h.i]
	h.i += 1
	if h.i == len(h.Colors) {
		h.i = 0
	}
	newInfo.Next = next
	newInfo.Color = next
	newMsg, err := json.Marshal(newInfo)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return newMsg
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
		Players: h.GameData.Players,
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
			close(client.Received)
			delete(h.Clients, client)
		}
	}
}
