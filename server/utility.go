package server

import (
	"encoding/json"
	"io"
	"math/rand"
)

//Treat These as constants. You can change COLORS though.
const CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func DecodeBody(data io.ReadCloser) (*ReqBody, error) {
	// Decode Json message from HTTP Request.
	// Send decoded into struct
	decoder := json.NewDecoder(data)
	decoder.DisallowUnknownFields()
	var body ReqBody
	for decoder.More() {
		err := decoder.Decode(&body)
		if err != nil {
			return nil, err
		}
	}
	return &body, nil
}

func MakeId() string {
	// Generates an id that unique amongst common room names
	id := ""
	for i := 0; i < 8; i++ {
		// All letters are valid. NO need to check
		id += string(CHARS[rand.Intn(len(CHARS))])
	}
	if IdExists(RoomStorage, id) {
		// This string has already been created
		return MakeId()
	}
	return id
}

func MakePin(room string) string {
	// Generates a password for others to connect to game
	pin := ""
	nums := "1234567890"
	for i := 0; i < 5; i++ {
		pin += string(nums[rand.Intn(len(nums))])
	}
	for _, val := range RoomStorage {
		if val.Hub.RoomData.Room == room && val.Hub.RoomData.Pin == pin {
			return MakePin(room)
		}
	}
	return pin
}
func IdExists(rooms Storage, id string) bool {
	// Checks if id exists in global games
	_, ok := rooms[id]
	return ok
}
