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
	// Addr := 76.192.124.46
	rand.Seed(time.Now().UnixNano())
	r := server.MakeRouter()
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
