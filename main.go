package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/BozeBro/ChainReaction/server"
)

func main() {
	// Makes it so that rand functions are actually pseudo random
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
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Println("Starting on :" + port)
	log.Fatal(srv.ListenAndServe())

}
