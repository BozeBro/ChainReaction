package server

import (
	"log"
	"net/http"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/mux"
)

// Struct to keep track of total gamers on the server
type PlayerCounter struct {
	Max     int
	Current int
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
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
					pc.Current--
					pc.Max = max(pc.Max, pc.Current)
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
	pc.Current++
	pc.Max = max(pc.Max, pc.Current)
}
