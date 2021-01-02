package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/BozeBro/ChainReaction/server"
)

func main() {
	// Makes it so that rand functions are actually pseudo random
	rand.Seed(time.Now().UnixNano())
	r := server.MakeRouter()
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
