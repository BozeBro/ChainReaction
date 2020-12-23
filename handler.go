package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	route := filepath.Join("website", "static", "html", "index.html")
	http.ServeFile(w, r, route)
}
func WaitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WE HERE")
	id := mux.Vars(r)["id"]
	if !IdExists(RoomStorage, id) {
		log.Println("CRYING")
		http.NotFound(w, r)
		return
	}
	//isleader := <-RoomStorage[id].Roles
	userRole := struct {
		Leader bool
	}{
		Leader: true,
	}
	route := "website/static/html/game.html"
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
		Roles:   make(chan bool, playerAmount)}
	func() { RoomStorage[id].Roles <- true }()
	http.Redirect(w, r, "/game/" + id, 302)
}
