package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Creates a route that will handle the routes for Chain Reaction
func MakeRouter() *mux.Router {
	// http.Dir uses directory of current working / dir where program started
	static := http.FileServer(http.Dir("./static"))
	roomStorage := make(Storage, 0)
	r := mux.NewRouter()
	r.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if _, ok := roomStorage[id]; ok {
			WSHandshake(roomStorage[id], w, r, roomStorage)
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/", HomeHandler)
	// Only the browser should be asking for the static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
	r.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		LobbyHandler(w, r, roomStorage)
	})
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/create/", func(w http.ResponseWriter, r *http.Request) { CreateHandler(w, r, roomStorage) }).Methods("POST")
	api.HandleFunc("/join/", func(w http.ResponseWriter, r *http.Request) { JoinHandler(w, r, roomStorage) }).Methods("POST")
	http.Handle("/", r)
	return r
}
