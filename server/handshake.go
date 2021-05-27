package server

import (
	"encoding/json"
	"log"
	"net/http"

	sock "github.com/BozeBro/ChainReaction/websocket"
	names "github.com/Pallinder/go-randomdata"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// upgrader upgrades a http connection to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Creates a websocket connection and starts the hub and client goroutines
// Might be changed to POST if a websocket server and http server are separated.
// Function couples the server and the websocket together so two servers are not needed
// Is in charge of stopping hub server
func WSHandshake(w http.ResponseWriter, r *http.Request, roomStorage Storage) {
	id := mux.Vars(r)["id"]
	hub := roomStorage[id]
	rolesws := hub.RoomData.Rolesws
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Person cannot use websockets
		log.Println(err)
		return
	}
	// Person didn't go through http route
	if len(rolesws) == 0 {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	isleader := <-rolesws
	username := <-hub.RoomData.Username
	// Only start hub server once.
	if isleader && !hub.Alive {
		go hub.Run()
		go func() {
			for {
				select {
				case <-hub.Stop:
					delete(roomStorage, id)
					return
				}
			}
		}()
	}
	// Arbitrarily large buffer to encourage make async programming
	// Prevent blocking when large amounts of animation data
	client := &sock.Client{
		Hub:      hub,
		Conn:     conn,
		Received: make(chan []byte, 1000),
		Leader:   isleader,
		Username: username,
	}
	hub.Register <- client
	go client.ReadMsg()
	go client.WriteMsg()
	if hub.RoomData.IsBot {
		botclient := &sock.Client{
			Hub:      hub,
			Received: make(chan []byte, 1000),
			Conn:     nil,
			Username: names.SillyName(),
		}
		hub.Register <- botclient
		go func() {
			move := botclient.Move()
			for {
				select {
				// dump out received msg to not fill chan
				case msg := <-botclient.Received:
					playInfo := new(sock.WSData)
					if err := json.Unmarshal(msg, playInfo); err != nil {
						log.Println(err)
						return
					}
					switch playInfo.Type {
					case "start", "move":
						canMove := playInfo.Turn == botclient.Color && len(hub.Colors) > 1
						if canMove {
							x, y := func(f string) (int, int) {
								if f == "rand" {
									return hub.Match.RandMove(botclient.Color)
								} else if f == "mm" {
									nextColor := ""
									for _, val := range hub.Colors {
										if val != botclient.Color {
											nextColor = val
										}
									}
									if nextColor == "" {
										log.Fatal("nextColor is nil: handshake line 109")
										return -1, -1
									}
									a := -100000
									b := 100000
									_, sq := hub.Match.Max(
										botclient.Color,
										nextColor,
										2,
										a,
										b,
										-6,
										-6,
									)
									return sq[0], sq[1]
								}

								return 1, 1
							}("mm")
							playInfo.X, playInfo.Y = x, y
							playInfo.Type = "move"
							if err := move(playInfo); err != nil {
								log.Println(err)
								return
							}
						}
					case "color":
						botclient.Color = playInfo.Color
					}
				case <-botclient.Stop:
					return
				}
			}
		}()
	}
}
