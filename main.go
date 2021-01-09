package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/BozeBro/ChainReaction/server"
)

// SUID | RW | RW | RW
const SUIDRWRR = 4664

func main() {
	// Makes it so that rand functions are actually pseudo random
	// Addr := 76.192.124.46
	rand.Seed(time.Now().UnixNano())
	finish := make(chan error)
	r, playerCounter := server.MakeRouter()
	ticker := time.NewTicker(time.Second * 24)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting on 127.0.0.1:8000")
	go func() {
		for {
			select {
			case <-ticker.C:
				var oldVal int
				file, err := ioutil.ReadFile("data.txt")
				if err != nil {
					err := ioutil.WriteFile("data.txt", []byte(fmt.Sprintf("%d", playerCounter.TotalPlayers)), SUIDRWRR)
					if err != nil {
						log.Println("Main.go:line 40: The file probably didn't exist.", err)
					}
					continue
				}
				if val := string(file); val == "" || val == " " {
					oldVal = 0
				} else {
					oldVal, err = strconv.Atoi(val)
					if err != nil {
						log.Printf("The value in the file was %d and the error was %s", oldVal, err)
						oldVal = 0
					}
				}
				if playerCounter.TotalPlayers > oldVal {
					err := ioutil.WriteFile("data.txt", []byte(fmt.Sprintf("%d", playerCounter.TotalPlayers)), SUIDRWRR)
					if err != nil {
						log.Printf("The error is occuring in the main file. Trying to write to a file. %s", err)
					}
				}
			case err := <-finish:
				log.Fatal("The server is dying down", err)
				return
			}
		}
	}()
	err := srv.ListenAndServe()
	finish <- err

}
