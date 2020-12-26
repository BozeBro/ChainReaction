package webserver

import (
	"encoding/json"
	"log"
)

type Hub struct {
	// alive ensures that only one hub will be created
	Alive bool
	// Mapping of clients. Unordered
	Clients map[*Client]bool

	// incoming broadcasting req from clients
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
	// Track who's turn it is
	I int
	// Way of Tracking who's turn it is
	Colors []string
}

type WSData struct {
	Type  string `json:"type"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
	Val   bool   `json:"val"`
	Next  string `json:"next"` // Next Color
}

func NewHub() *Hub {
	return &Hub{
		Alive:      false,
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
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
	defer func() { h.Alive = false }()
	// wait for register, unregister, or broadcast chan to be filled
	for {
		select {
		case client := <-h.Register:
			client.Color = h.GetUniqueColor(RandomColor())
			var wsData = &WSData{Color: client.Color, Type: "color"}
			log.Println(client.Color)
			payload, err := json.Marshal(wsData)
			if err != nil {
				// This should never happen.
				log.Fatal(err)
				return
			}
			h.Clients[client] = true
			client.Received <- payload
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Received)
			}
		case message := <-h.Broadcast:
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
		// We will only see "start" in beginning of each game
		h.Colors = make([]string, len(h.Clients))
		i := 0
		for k, _ := range h.Clients {
			h.Colors[i] = k.Color
			i++
		}
	}
	next := h.Colors[h.I]
	h.I++
	if h.I >= len(h.Colors) {
		h.I = 0
	}
	newInfo.Next = next
	newInfo.Color = next
	newMsg, err := json.Marshal(newInfo)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println(next, h.Colors)
	return newMsg
}
