package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func LobbyHandler(w http.ResponseWriter, r *http.Request, roomStorage Storage) {
	// Handler that serves game.file.
	// Initialized to show waiting screen
	// Create and Join Handler will route everything here
	id := mux.Vars(r)["id"]
	if !IdExists(roomStorage, id) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	// A person is joining via URL directly and not GUI PIN/NAME system
	if len(roomStorage[id].Roles) == 0 && roomStorage[id].Hub.Alive {
		// Leader is in the game and it is not full
		notFull := roomStorage[id].Hub.RoomData.Players+1 <= roomStorage[id].Hub.RoomData.Max
		if notFull {
			roomStorage[id].Roles <- false
			roomStorage[id].Rolesws <- false
			http.Redirect(w, r, "/game/"+id, http.StatusFound)
			return
		}
		// Can't play in a game if capacity is reached.
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if len(roomStorage[id].Roles) == 0 && !roomStorage[id].Hub.Alive {
		// There is no game / leader failed to connect to websocket
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	isleader := <-roomStorage[id].Roles
	userRole := struct {
		Leader  bool
		Players int
		Max     int
		Pin     string
		Room    string
	}{
		Leader:  isleader,
		Players: roomStorage[id].Hub.RoomData.Players,
		Max:     roomStorage[id].Hub.RoomData.Max,
		Room:    roomStorage[id].Hub.RoomData.Room,
		Pin:     roomStorage[id].Hub.RoomData.Pin,
	}
	route := "static/html/game.html"
	gameFile := template.Must(template.ParseFiles(route))
	if err := gameFile.Execute(w, userRole); err != nil {
		log.Println(err)
		return
	}
}
