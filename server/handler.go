package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	sock "github.com/BozeBro/ChainReaction/websocket"
)

// Storage is the type that stores all Games
// id is the key with information about it as a value
type Storage map[string]*sock.Hub

// HomeHandler send the index.html page at root path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	route := filepath.Join("static", "html", "index.html")
	http.ServeFile(w, r, route)
}

// JoinHandler handles POST requests to join a game.
// Redirects users back to root with a message if there is an error.
// Sends back a path with an id to go to to if there is no error.
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
	// rooms contains ids of active rooms with the same name as POST information
	rooms := make([]string, 0)
	for id, hub := range roomStorage {
		sameRoom := body.Room == hub.RoomData.Room
		if sameRoom {
			rooms = append(rooms, id)
		}
	}
	if len(rooms) == 0 {
		http.Error(w, "Room doesn't exist", http.StatusNotFound)
		return
	}
	for _, id := range rooms {
		hub := roomStorage[id]
		samePin := body.Pin == hub.RoomData.Pin
		notFull := hub.RoomData.Players+1 <= hub.RoomData.Max
		if samePin && !notFull {
			http.Error(w, "The room is full", http.StatusForbidden)
			return
		}
		if samePin && notFull {
			hub.RoomData.Roles <- false
			hub.RoomData.Rolesws <- false
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(id))
			return
		}
	}
	http.Error(w, "Wrong Pin", http.StatusNotAcceptable)
}

// CreateHandler creates a room.
// Sends back a path with an id to go to along with leader permissions.
// Redirects user back to root path with a message if there is an error / problems.
// BUG: CreateHandler doesn't handle rooms with special character names
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
	if body.Players == "" || body.Room == "" {
		http.Error(w, "Empty values", http.StatusConflict)
		return
	}
	body.Room = strings.ReplaceAll(body.Room, " ", "")
	playerAmount, err := strconv.Atoi(body.Players)
	// user POST value that is not a number
	if err != nil {
		log.Printf("Someone sent a strange playerAmount of %v", playerAmount)
		http.Error(w, "Nice Try nerd", http.StatusBadRequest)
		return
	}
	// Context-like values being created
	id := MakeId(roomStorage)
	pin := MakePin(body.Room, roomStorage)
	gameinfo := &sock.RoomData{
		Room:    body.Room,
		Pin:     pin,
		Max:     playerAmount,
		Players: 0, // Correct players will be in joinHandler
	}
	hub := sock.NewHub(gameinfo)
	hub.RoomData.Roles = make(chan bool, playerAmount)
	hub.RoomData.Rolesws = make(chan bool, playerAmount)
	roomStorage[id] = hub
	roomStorage[id].RoomData.Roles <- true
	roomStorage[id].RoomData.Rolesws <- true
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(id))
}
