package main

import (
	"log"
	"net/http"
)

type WSData struct {
	Type  string `json:"type"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
	Val   bool  `json:"val"`
	Next string `json:"next"` // Next Color
}

func WSHandshake(g *GameData, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	/*
	if len(g.Rolesws) == 0 {
		http.Error(w, "You did not enter properly", 409)
		return
	}
	*/
	//isleader := <-g.Rolesws
	isleader := true
	go func() {
		if isleader {
			g.Hub.Run()
		}
	}()
	client := &Client{hub: g.Hub, 
		conn: conn, 
		received: make(chan []byte, 256), 
		Leader: isleader,}
	client.hub.register <- client
	go client.readMsg()
	go client.writeMsg()
}
