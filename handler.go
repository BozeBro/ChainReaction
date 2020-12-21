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
	params := mux.Vars(r)
	id := params["id"]
	log.Println(id)
	http.ServeFile(w, r, "./website/static/html/waiting.html")
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
		http.Error(w, "You didn't send an integer value for player amount", 400)
		log.Fatal(err)
	}
	// Create Proper Unique Data
	id := MakeId()
	pin := MakePin(body.Room)
	RoomStorage[id] = &GameData{Players: playerAmount, Room: body.Room, Pin: pin}
	http.Redirect(w, r, "/game/"+id, 303)
	// Websocket Connection that monitors players in the room.
	// Once Mod/Admin presses "Start", then begin the Game
}
