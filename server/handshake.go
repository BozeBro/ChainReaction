package server

import (
	"log"
	"net/http"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/mux"
)

func WSHandshake(g *GameData, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if len(g.Rolesws) == 0 {
		http.Error(w, "You did not enter properly", http.StatusConflict)
		return
	}
	if len(g.Rolesws) == 0 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
	isleader := <-g.Rolesws
	if isleader && !g.Hub.Alive {
		go g.Hub.Run()
		go func() {
			for {
				select {
				case <-g.Hub.Stop:
					id := mux.Vars(r)["id"]
					delete(RoomStorage, id)
					return
				case <-g.Hub.Delete:
					g.Hub.RoomData.Players -= 1
					g.Hub.Update()
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
	client.Hub.Register <- client
	go client.ReadMsg()
	go client.WriteMsg()
}
