package server

import (
	"log"
	"net/http"

	sock "github.com/BozeBro/ChainReaction/websocket"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// upgrader upgrades a http connection to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Creates a websocket connection and starts the hub and client goroutines
// Might be changed to POST if a websocket server and http server are separated.
// Function couples the server and the websocket together so two servers are not needed
// Is in charge of stopping hub server
func WSHandshake(w http.ResponseWriter, r *http.Request, roomStorage Storage) {
	id := mux.Vars(r)["id"]
	hub := roomStorage[id]
	rolesws := hub.RoomData.Rolesws
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Person cannot use websockets
		log.Println(err)
		return
	}
	// Person didn't go through http route
	if len(rolesws) == 0 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	isleader := <-rolesws
	// Only start hub server once.
	if isleader && !hub.Alive {
		go hub.Run()
		go func() {
			for {
				select {
				case <-hub.Stop:
					delete(roomStorage, id)
					return
				}
			}
		}()
	}
	// Arbitrarily large buffer to encourage make async programming
	// Prevent blocking when large amounts of animation data
	client := &sock.Client{
		Hub:      hub,
		Conn:     conn,
		Received: make(chan []byte, 1000),
		Leader:   isleader,
	}
	hub.Register <- client
	go client.ReadMsg()
	go client.WriteMsg()
}
