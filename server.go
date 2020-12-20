package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/mux"
)
type Storage map[string]*GameData

var RoomStorage = make(Storage, 0)

type GameData struct {
	Room, Pin string
	Players string
}
type ReqBody struct {
	Pin, Room, Players, Name string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/static/html/index.html")
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
	_, err = strconv.Atoi(body.Players)
	if err != nil {
		// Someone attempting some hacks
		http.Error(w, "You didn't send a numerical value for player amount", 400)
		log.Fatal(err)
	}
	// Create Proper Unique Data
	id := MakeId()
	pin := MakePin()
	log.Println(RoomStorage)
	RoomStorage[id] = &GameData{Players: body.Players, Room: body.Room, Pin: pin}
	// Redirect Person to Game Room. /game/{Identity} URL should have a unique identifier
	// Websocket Connection that monitors players in the room.
	// Once Mod/Admin presses "Start", then begin the Game
	http.Error(w, "Not Implemented", 500)
}
func main() {
	static := http.FileServer(http.Dir("./website/static"))
	r := mux.NewRouter()
	hub := newHub()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if !hub.alive {
			hub.alive = true
			go hub.run()
		}
		wSHandshake(hub, w, r)
	})
	r.HandleFunc("/", HomeHandler).
		Methods("GET")
	r.PathPrefix("/css/{file}").Handler(static)
	r.PathPrefix("/js/{file}").Handler(static)
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/create", CreateHandler)
	api.HandleFunc("/join", JoinHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("serving at 127.0.0.1:8000")
	log.Fatal(srv.ListenAndServe())
}
