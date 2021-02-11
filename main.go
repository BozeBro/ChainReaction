package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/BozeBro/ChainReaction/server"
	"github.com/joho/godotenv"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// Multiplexer that chain reaction runs on
	r := server.MakeRouter()
	// dynamic port assigned by Heroku
	port := os.Getenv("PORT")
	if port == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("$PORT must be set")
		} else {
			// running locally
			port = os.Getenv("PORT")
		}
	}
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
	log.Println("Starting on :" + port)
	log.Fatal(srv.ListenAndServe())

}
