package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler that serves game file
// Create and Join Handler will route here
// Redirects people who reach here by URL back to Join to be stored in context
func LobbyHandler(w http.ResponseWriter, r *http.Request, roomStorage Storage) {
	id := mux.Vars(r)["id"]
	if !IdExists(roomStorage, id) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	hub := roomStorage[id]
	roomData := hub.RoomData
	// A person is joining via URL directly and not GUI PIN/NAME system
	if len(roomData.Roles) == 0 && hub.Alive {
		// Leader is in the game and it is not full
		notFull := roomData.Players+1 <= roomData.Max
		if notFull {
			roomData.Roles <- false
			roomData.Rolesws <- false
			http.Redirect(w, r, "/game/"+id, http.StatusFound)
			return
		}
		// Can't play in a game if capacity is reached.
		http.Redirect(w, r, "/", http.StatusFound)
		return
		// There is no game / leader failed to connect to websocket
	} else if len(roomData.Roles) == 0 && !hub.Alive {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	isleader := <-roomData.Roles
	userRole := struct {
		Leader  bool
		Players int
		Max     int
		Pin     string
		Room    string
	}{
		Leader:  isleader,
		Players: roomData.Players,
		Max:     roomData.Max,
		Room:    roomData.Room,
		Pin:     roomData.Pin,
	}
	route := "static/html/game.html"
	gameFile := template.Must(template.ParseFiles(route))
	if err := gameFile.Execute(w, userRole); err != nil {
		log.Println(err)
		return
	}
}
