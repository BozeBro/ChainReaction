package main

import (
	"log"
)

type Hub struct {
	alive bool
	// Mapping of clients
	clients map[*Client]bool

	// incoming broadcasting req from clients
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		alive:      false,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	defer func() { h.alive = false }()
	// wait for register, unregister, or broadcast chan to be filled
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.received)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.received <- message:
				default:
					close(client.received)
					delete(h.clients, client)
				}
			}
		}
	}
}
