package server

import (
	"net/http"

	"github.com/BozeBro/ChainReaction/webserver"
	"github.com/gorilla/mux"
)

// The string will be the id
type Storage map[string]*GameData

var RoomStorage = make(Storage, 0)

type GameData struct {
	Room, Pin string
	Players   int
	Hub       *webserver.Hub // The game server
	Roles     chan bool // Send roles to handler
	Rolesws   chan bool // send roles to handler of websockets
	Max int // Maximum amount of players
}
type ReqBody struct {
	Pin, Room, Players, Name string
}

func MakeRouter() *mux.Router {
	// http.Dir uses directory of current working / dir where program started
	static := http.FileServer(http.Dir("./static"))
	r := mux.NewRouter()
	//hub := newHub()
	r.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		WSHandshake(RoomStorage[id], w, r)
	})
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
	r.HandleFunc("/game/{id}", WaitHandler).Methods("GET")
	r.HandleFunc("/game/{id}/join", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/create/", CreateHandler).Methods("POST")
	api.HandleFunc("/join/", JoinHandler).Methods("POST")
	http.Handle("/", r)
	return r
}
