package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/static/html/index.html")
}
func WaitHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if !IdExists(RoomStorage, id) {
		http.NotFound(w, r)
		return
	}
	isleader := <-RoomStorage[id].Roles
	if isleader {
		http.ServeFile(w, r, "./website/static/html/waitingLeader.html")
	} else {
		http.ServeFile(w, r, "./website/static/html/waitingPlayer.html")
	}
	// Websocket Connection that monitors players in the room.
	// Once Mod/Admin presses "Start", then begin the Game
}
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(body)
	// Check if Room Exists
	// Join that Room
}
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := DecodeBody(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Take in anything approach.
	if body.Players == "" || body.Room == "" {
		http.Error(w, "Empty values", 409)
		return
	}
	isUnique := func() bool {
		for _, val := range RoomStorage {
			if (*val).Room == body.Room {
				return false
			}
		}
		return true
	}()
	if !isUnique {
		log.Println(RoomStorage)
		http.Error(w, "The room is already taken. Try another name", 409)
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
	go func() { RoomStorage[id].Roles <- true }()
	http.Redirect(w, r, "/game/"+id, 303)

}
