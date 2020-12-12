package main

import (
	"log"
	"net/http"

)


func wSHandshake(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: h, conn: conn, received: make(chan []byte, 256)}
	client.hub.register <- client
	go client.readMsg()
	go client.writeMsg()
}
