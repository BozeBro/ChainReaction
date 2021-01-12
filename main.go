package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/BozeBro/ChainReaction/server"
)

// SUID | RW | RW | RW
const SUIDRWRR = 4664

func main() {
	// Makes it so that rand functions are actually pseudo random
	// Addr := 76.192.124.46
	rand.Seed(time.Now().UnixNano())
	r := server.MakeRouter()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting on :" + port)
	log.Fatal(srv.ListenAndServe())

}
