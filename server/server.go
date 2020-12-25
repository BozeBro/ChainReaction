package server

import (
	"log"
	"net/http"
	"time"
	"math/rand"
	"github.com/gorilla/mux"
)

// The string will be the id
type Storage map[string]*GameData

var RoomStorage = make(Storage, 0)

type GameData struct {
	Room, Pin string
	Players   int
	Hub       *Hub
	Roles     chan bool
	Rolesws   chan bool
}
type ReqBody struct {
	Pin, Room, Players, Name string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	static := http.FileServer(http.Dir("./website/static"))
	r := mux.NewRouter()
	//hub := newHub()
	r.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		WSHandshake(RoomStorage[id], w, r)
	})
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.PathPrefix("/website/").Handler(http.StripPrefix("/website/static/", static))
	r.HandleFunc("/game/{id}", WaitHandler).Methods("GET")
	r.HandleFunc("/game/{id}/join", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/create", CreateHandler).Methods("POST")
	api.HandleFunc("/join", JoinHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("serving at 127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
