package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/BozeBro/ChainReaction/webserver"
)

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
	if len(RoomStorage[id].Roles) == 0 {
		log.Println("Hijacking?")
		return
	}
	//isleader := <-RoomStorage[id].Roles
	isleader := true
	userRole := struct {
		Leader bool
	}{
		Leader: isleader,
	}
	route := "static/html/game.html"
	gameFile := template.Must(template.ParseFiles(route))
	if err := gameFile.Execute(w, userRole); err != nil {
		log.Fatal(err)
		return
	}
}
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(body)
	// Check if Room Exists
	// Join that Room
}
func CreateHandler(w http.ResponseWriter, r *http.Request) {
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
	RoomStorage[id] = &GameData{
		Players: playerAmount,
		Room:    body.Room,
		Pin:     pin,
		Hub:     webserver.NewHub(),
		Roles:   make(chan bool, playerAmount),
		Rolesws: make(chan bool, playerAmount),
	}
	func() {
		RoomStorage[id].Roles <- true
		RoomStorage[id].Rolesws <- true
	}()
	http.Redirect(w, r, "/game/"+id+"/join", 302)
}
func WSHandshake(g *GameData, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	/*
		if len(g.Rolesws) == 0 {
			http.Error(w, "You did not enter properly", 409)
			return
		}
	*/
	//isleader := <-g.Rolesws
	isleader := true
	go func() {
		if isleader && !g.Hub.Alive {
			g.Hub.Alive = true
			g.Hub.Run()
		}
	}()
	client := &webserver.Client{
		Hub: g.Hub,
		Conn:     conn,
		Received: make(chan []byte, 256),
		Leader:   isleader}
	client.Hub.Register <- client
	go client.ReadMsg()
	go client.WriteMsg()
}
