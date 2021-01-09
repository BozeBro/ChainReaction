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
	id := mux.Vars(r)["id"]
	if !IdExists(roomStorage, id) {
		log.Println("Room Doesn't Exist")
		http.NotFound(w, r)
		return
	}

	if len(roomStorage[id].Roles) == 0 && roomStorage[id].Hub.Alive {
		notFull := roomStorage[id].Hub.RoomData.Players+1 <= roomStorage[id].Hub.RoomData.Max
		if notFull {
			go func() {
				roomStorage[id].Roles <- false
				roomStorage[id].Rolesws <- false
			}()
			http.Redirect(w, r, "/game/"+id, http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if len(roomStorage[id].Roles) == 0 && !roomStorage[id].Hub.Alive {
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
