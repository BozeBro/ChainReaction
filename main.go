package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type Layout struct {
	/*
	The template data applied to every html file
	*/
	Data []string
	Head string
}

func FileServer(w http.ResponseWriter, r *http.Request) {
	/*
	This will handle all the html files coming at "/"
	gorilla/mux handles nonexistant pages with a 404
	*/
	filepath := ""
	switch r.URL.Path {
	case "/", "/index", "/home":
		filepath = "./static/index.html"
	case "/newbie":
		filepath = "./static/new.html"
	}
	func() {
		http.ServeFile(w, r, filepath)
	}()
}
func main() {
	var layout string

	flag.StringVar(&layout, "layout", "./website/static/layout.html", "Location of the layout.html file path")
	flag.Parse()

	tmpl := template.Must(template.ParseFiles(layout))
	r := mux.NewRouter()
	r.HandleFunc("/", FileServer)
	r.HandleFunc("/index", FileServer)
	r.HandleFunc("/home", FileServer)
	r.HandleFunc("/newbie", FileServer)
	r.HandleFunc("/layout", func(w http.ResponseWriter, r *http.Request) {
		data := Layout{
			Data: []string{
				"maybe",
				"someday",
				"I",
				"will",
				"swim",
			},
			Head: "MY HEAD",
		}
		log.Fatal(tmpl.Execute(w, data))

	})

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("STARTING port 8000")
	log.Fatal(srv.ListenAndServe())
}
