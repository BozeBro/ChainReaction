package server

import (
	"log"
	"net/http"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/mux"
)

// Struct to keep track of total gamers on the server
type PlayerCounter struct {
	TotalPlayers int
}

func WSHandshake(g *GameData, w http.ResponseWriter, r *http.Request, roomStorage Storage, pc *PlayerCounter) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if len(g.Rolesws) == 0 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	isleader := <-g.Rolesws
	if isleader && !g.Hub.Alive {
		go g.Hub.Run()
		go func() {
			for {
				select {
				case <-g.Hub.Stop:
					g.Hub.CloseChans()
					delete(roomStorage, mux.Vars(r)["id"])
					return
				case <-g.Hub.Leaver:
					pc.TotalPlayers--
				}
			}
		}()
	}
	client := &webserver.Client{
		Hub:      g.Hub,
		Conn:     conn,
		Received: make(chan []byte, 256),
		Leader:   isleader,
	}
	go func() {
		client.Hub.Register <- client
	}()
	go client.ReadMsg()
	go client.WriteMsg()
	pc.TotalPlayers++
}
