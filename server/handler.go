package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/websocket"
)

/*
The main Handlers that are found here are the Join and Create Handlers. To see the lobby handler, see lobby.go
*/
// The string will be the id
type Storage map[string]*GameData

type GameData struct {
	Hub     *webserver.Hub // The game server
	Roles   chan bool      // Send roles to handler
	Rolesws chan bool      // send roles to handler of websockets
}

// Global variable that allows us to upgrade a connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	route := filepath.Join("static", "html", "index.html")
	http.ServeFile(w, r, route)
}
func JoinHandler(w http.ResponseWriter, r *http.Request, roomStorage Storage) {
	if r.Method != "POST" {
		log.Print("http Method was illegal for Join")
		return
	}
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	// Do multiple loops so we know what to tell the user
	rooms := make([]string, 0)
	for id, data := range roomStorage {
		sameRoom := body.Room == data.Hub.RoomData.Room
		if sameRoom {
			rooms = append(rooms, id)
		}
	}
	if len(rooms) == 0 {
		http.Error(w, "Room doesn't exist", http.StatusNotFound)
		return
	}
	for _, id := range rooms {
		// Checks if room has vacancy
		data := roomStorage[id]
		samePin := body.Pin == data.Hub.RoomData.Pin
		notFull := data.Hub.RoomData.Players+1 <= data.Hub.RoomData.Max
		if samePin && !notFull {
			http.Error(w, "The room is full", http.StatusForbidden)
			return
		}
		if samePin && notFull {
			go func() {
				data.Roles <- false
				data.Rolesws <- false
			}()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(id))
			return
		}
	}
	http.Error(w, "Wrong Pin", http.StatusNotAcceptable)
}

//Creates room. Redirects to empty handler. Redirects to JoinHandler
func CreateHandler(w http.ResponseWriter, r *http.Request, roomStorage Storage) {
	if r.Method != "POST" {
		log.Print("http Method was illegal for Create")
		return
	}
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	// Take in anything approach.
	if body.Players == "" || body.Room == "" {
		http.Error(w, "Empty values", http.StatusConflict)
		return
	}
	playerAmount, err := strconv.Atoi(body.Players)
	if err != nil {
		// Someone attempting some hacks
		// Might be unneccessary though
		log.Println(err, "Player count was not a number?? shouldn't be possible")
		http.Error(w, "Nice Try nerd", http.StatusBadRequest)
		return
	}
	// Create Proper Unique Data
	id := MakeId(roomStorage)
	pin := MakePin(body.Room, roomStorage)
	gameinfo := &webserver.RoomData{
		Room:    body.Room,
		Pin:     pin,
		Max:     playerAmount,
		Players: 0, // Correct players will be in joinHandler
	}
	roomStorage[id] = &GameData{
		Hub:     webserver.NewHub(gameinfo),
		Roles:   make(chan bool, playerAmount),
		Rolesws: make(chan bool, playerAmount),
	}
	go func() {
		roomStorage[id].Roles <- true
	}()
	go func() {
		roomStorage[id].Rolesws <- true
	}()
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.Write([]byte(id))
	userRole := struct {
		Leader  bool
		Players int
		Max     int
		Pin     string
		Room    string
	}{
		Leader:  true,
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
