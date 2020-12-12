package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var dir string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()
	hub := newHub()
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir(dir + "/static")))
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if !hub.alive {
			hub.alive = true
			go hub.run()
		}
		wSHandshake(hub, w, r)
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(dir)))
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
