package websocket

import (
	"encoding/json"
	"log"
)

func (bot *Client) HandleMsg(msg []byte) {
	playInfo := &WSData{}
	if err := json.Unmarshal(msg, playInfo); err != nil {
		log.Println(err)
		return
	}
	bot.PlayMove(playInfo)
}

func (bot *Client) PlayMove(playInfo *WSData) {
	hub := bot.Hub
	canMove := playInfo.Turn == bot.Color && len(hub.Colors) > 1
	if !canMove {
		return
	}
	move := bot.Move()
	// Game Options. Editable only by programmer
	options := "mm"
	playInfo.Type = "move"
	if options != "mm" {
		playInfo.X, playInfo.Y = hub.Match.RandMove(bot.Color)
		if err := move(playInfo); err != nil {
			log.Fatal(err)
		}
		return
	}
	nextColor := ""
	for _, val := range hub.Colors {
		if val != bot.Color {
			nextColor = val
		}
	}
	if nextColor == "" {
		log.Fatal("nextColor is nil: msgresbot.go line 38")
	}
	a, b := -100000, 100000
	depth := 3
	index := bot.Hub.i
	// movedx, movedy is a negative number to symbolize it is not used.
	_, sq := hub.Match.Max(
		bot.Color,
		nextColor,
		depth,
		a,
		b,
		-1,
		-1,
	)
	bot.Hub.i = index
	playInfo.X, playInfo.Y = sq[0], sq[1]
	if err := move(playInfo); err != nil {
		log.Println(err)
		return
	}

}
