package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// MakeRouter creates a router that will handle the routes for Chain Reaction.
func MakeRouter() *mux.Router {
	// static handles all front end files.
	static := http.FileServer(http.Dir("./static"))
	// roomStorage tracks all active games
	roomStorage := make(Storage, 0)
	r := mux.NewRouter()
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "favicon.ico") })
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "robots.txt") })
	r.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if _, ok := roomStorage[id]; ok {
			WSHandshake(w, r, roomStorage)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/", HomeHandler)
	// Only the browser should be asking for the static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
	r.HandleFunc("/game/{id:[a-zA-Z0-9]{8}}", func(w http.ResponseWriter, r *http.Request) {
		LobbyHandler(w, r, roomStorage)
	})
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/create/", func(w http.ResponseWriter, r *http.Request) { CreateHandler(w, r, roomStorage) }).Methods("POST")
	api.HandleFunc("/join/", func(w http.ResponseWriter, r *http.Request) { JoinHandler(w, r, roomStorage) }).Methods("POST")
	http.Handle("/", r)
	return r
}
