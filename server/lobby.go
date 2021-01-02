package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func LobbyHandler(w http.ResponseWriter, r *http.Request) {
	// Handler that serves game.file.
	// Initialized to show waiting screen
	id := mux.Vars(r)["id"]
	if !IdExists(RoomStorage, id) {
		log.Println("Room Doesn't Exist")
		http.NotFound(w, r)
		return
	}

	if len(RoomStorage[id].Roles) == 0 && RoomStorage[id].Hub.Alive {
		notFull := RoomStorage[id].Hub.RoomData.Players+1 <= RoomStorage[id].Hub.RoomData.Max
		if notFull {
			go func() {
				RoomStorage[id].Roles <- false
				RoomStorage[id].Rolesws <- false
			}()
			http.Redirect(w, r, "/game/"+id, http.StatusFound)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if len(RoomStorage[id].Roles) == 0 && !RoomStorage[id].Hub.Alive {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	isleader := <-RoomStorage[id].Roles
	userRole := struct {
		Leader  bool
		Players int
		Max     int
		Pin     string
		Room    string
	}{
		Leader:  isleader,
		Players: RoomStorage[id].Hub.RoomData.Players,
		Max:     RoomStorage[id].Hub.RoomData.Max,
		Room:    RoomStorage[id].Hub.RoomData.Room,
		Pin:     RoomStorage[id].Hub.RoomData.Pin,
	}
	route := "static/html/game.html"
	gameFile := template.Must(template.ParseFiles(route))
	if err := gameFile.Execute(w, userRole); err != nil {
		log.Fatal(err)
		return
	}
}
