package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/websocket"
)

// The string will be the id
type Storage map[string]*GameData

type GameData struct {
	Hub     *webserver.Hub // The game server
	Roles   chan bool      // Send roles to handler
	Rolesws chan bool      // send roles to handler of websockets
}

var RoomStorage = make(Storage, 0)

// Global variable that allows us to upgrade a connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	route := filepath.Join("static", "html", "index.html")
	http.ServeFile(w, r, route)
}
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Do multiple loops so we know what to tell the user
	rooms := make([]string, 0)
	for id, data := range RoomStorage {
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
		data := RoomStorage[id]
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
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
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
		http.Error(w, "Nice Try nerd", http.StatusBadRequest)
		log.Println(err)
		return
	}
	// Create Proper Unique Data
	id := MakeId()
	pin := MakePin(body.Room)
	gameinfo := &webserver.RoomData{
		Room:    body.Room,
		Pin:     pin,
		Max:     playerAmount,
		Players: 0, // Correct players will be in joinHandler
	}
	RoomStorage[id] = &GameData{
		Hub:     webserver.NewHub(gameinfo),
		Roles:   make(chan bool, playerAmount),
		Rolesws: make(chan bool, playerAmount),
	}
	func() {
		RoomStorage[id].Roles <- true
		RoomStorage[id].Rolesws <- true
	}()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(id))
}
