package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ReqBody struct {
	Pin, Room, Players, Name string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./website/static/html/index.html")
}
func DecodeBody(data io.ReadCloser) (*ReqBody, error) {
	decoder := json.NewDecoder(data)
	var body ReqBody
	for decoder.More() {
		err := decoder.Decode(&body)
		if err != nil {
			return nil, err
		}
	}
	return &body, nil
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
	log.Println(body)
	// Create Room
	// Add room to created Rooms
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
