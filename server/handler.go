package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/mux"
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
func WaitHandler(w http.ResponseWriter, r *http.Request) {
	// Handler that serves game.file.
	// Initialized to show waiting screen
	id := mux.Vars(r)["id"]
	if !IdExists(RoomStorage, id) {
		log.Println("Room Doesn't Exist")
		http.NotFound(w, r)
		return
	}

	if len(RoomStorage[id].Roles) == 0 && RoomStorage[id].Hub.Alive {
		notFull := RoomStorage[id].Hub.GameData.Players+1 <= RoomStorage[id].Hub.GameData.Max
		if notFull {
			go func() {
				RoomStorage[id].Roles <- false
				RoomStorage[id].Rolesws <- false
			}()
			http.Redirect(w, r, "/game/"+id, 302)
			return
		}
		http.Redirect(w, r, "/", 302)
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
		Players: RoomStorage[id].Hub.GameData.Players,
		Max:     RoomStorage[id].Hub.GameData.Max,
		Room:    RoomStorage[id].Hub.GameData.Room,
		Pin:     RoomStorage[id].Hub.GameData.Pin,
	}
	route := "static/html/game.html"
	gameFile := template.Must(template.ParseFiles(route))
	if err := gameFile.Execute(w, userRole); err != nil {
		log.Fatal(err)
		return
	}
}
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(RoomStorage)
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Do multiple loops so we know what to tell the user
	rooms := make([]string, 0)
	for id, gameData := range RoomStorage {
		sameRoom := body.Room == gameData.Hub.GameData.Room
		log.Println(sameRoom)
		log.Println(body.Room, gameData.Hub.GameData.Room)
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
		gameData := RoomStorage[id]
		samePin := body.Pin == gameData.Hub.GameData.Pin
		notFull := gameData.Hub.GameData.Players+1 <= gameData.Hub.GameData.Max
		if samePin && !notFull {
			http.Error(w, "The room is full", http.StatusForbidden)
			return
		}
		if samePin && notFull {
			go func() {
				gameData.Roles <- false
				gameData.Rolesws <- false
			}()
			http.Redirect(w, r, "/game/"+id+"/join", 302)
			return
		}
	}
	http.Error(w, "Wrong Pin", 406)
}
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Creates room. Redirects to empty handler. Redirects to joinHandler
	*/
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Take in anything approach.
	if body.Players == "" || body.Room == "" {
		http.Error(w, "Empty values", 409)
		return
	}
	playerAmount, err := strconv.Atoi(body.Players)
	if err != nil {
		// Someone attempting some hacks
		// Might be unneccessary though
		http.Error(w, "Nice Try nerd", 400)
		log.Println(err)
		return
	}
	// Create Proper Unique Data
	id := MakeId()
	pin := MakePin(body.Room)
	gameData := &webserver.GameData{
		Room:    body.Room,
		Pin:     pin,
		Max:     playerAmount,
		Players: 0, // Correct players will be in joinHandler
	}
	RoomStorage[id] = &GameData{
		Hub:     webserver.NewHub(gameData),
		Roles:   make(chan bool, playerAmount),
		Rolesws: make(chan bool, playerAmount),
	}
	func() {
		RoomStorage[id].Roles <- true
		RoomStorage[id].Rolesws <- true
	}()
	http.Redirect(w, r, "/game/"+id+"/join", 302)
}
