package main

import (
	"encoding/json"
	"log"
)

type Hub struct {
	// alive ensures that only one hub will be created
	alive bool
	// Mapping of clients. Unordered
	Clients map[*Client]bool

	// incoming broadcasting req from clients
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	// Track who's turn it is
	I int
	// Way of Tracking who's turn it is
	Colors []string
}

func newHub() *Hub {
	return &Hub{
		alive:      false,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
	defer func() { h.alive = false }()
	// wait for register, unregister, or broadcast chan to be filled
	for {
		select {
		case client := <-h.register:
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
			client.received <- payload
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.received)
			}
		case message := <-h.broadcast:
			newMsg := h.EditMsg(message)
			for client := range h.Clients {
				select {
				case client.received <- newMsg:
				default:
					close(client.received)
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
			log.Print(k)
			log.Println(" Client is")
			h.Colors[i] = k.Color
			i++
		}
	}
	log.Print("h.I is")
	log.Println(h.I)
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
