package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ReqBody struct {
	Pin, Room, Players, Name string
}

func MakeRouter() *mux.Router {
	// http.Dir uses directory of current working / dir where program started
	static := http.FileServer(http.Dir("./static"))
	r := mux.NewRouter()
	r.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if _, ok := RoomStorage[id]; ok {
			WSHandshake(RoomStorage[id], w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/", HomeHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
	r.HandleFunc("/game/{id}", WaitHandler)
	r.HandleFunc("/game/{id}/join", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/create/", CreateHandler).Methods("POST")
	api.HandleFunc("/join/", JoinHandler).Methods("POST")
	http.Handle("/", r)
	return r
}
