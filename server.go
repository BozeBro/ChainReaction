package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// The string will be the id
type Storage map[string]*GameData

var RoomStorage = make(Storage, 0)

type GameData struct {
	Room, Pin string
	Players   int
	Roles chan bool
}
type ReqBody struct {
	Pin, Room, Players, Name string
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
	r.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		go func ()  {
			id := mux.Vars(r)["id"]
			RoomStorage[id].Roles <- false
		}()
		WaitHandler(w, r)
	})
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
