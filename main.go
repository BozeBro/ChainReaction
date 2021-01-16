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
	rand.Seed(time.Now().UnixNano())
	// Multiplexer that chain reaction runs on
	r := server.MakeRouter()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// long
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
	log.Println("Starting on :" + port)
	log.Fatal(srv.ListenAndServe())

}
