package server

import (
	"encoding/json"
	"io"
	"math/rand"
)

// All the chars that MakePin can utilize
const CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// Valid request information that a user can send
// It is stored in the context
type ReqBody struct {
	Pin, Room, Players, Name string
}

// Decode Json message from HTTP Request.
// Prevents strange values.
func DecodeBody(data io.ReadCloser) (*ReqBody, error) {
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

// Generates a unique id containing only numbers
func MakeId(roomStorage Storage) string {
	id := ""
	for i := 0; i < 8; i++ {
		// All letters are valid. NO need to check
		id += string(CHARS[rand.Intn(len(CHARS))])
	}
	// This string has already been created
	if IdExists(roomStorage, id) {
		return MakeId(roomStorage)
	}
	return id
}

// MakePin Generates a unique PIN with length 5.
func MakePin(room string, roomStorage Storage) string {
	pin := ""
	nums := "1234567890"
	for i := 0; i < 5; i++ {
		pin += string(nums[rand.Intn(len(nums))])
	}
	for _, val := range roomStorage {
		if val.Hub.RoomData.Room == room && val.Hub.RoomData.Pin == pin {
			return MakePin(room, roomStorage)
		}
	}
	return pin
}

// IdExists checks if id exists in Storage variable
func IdExists(rooms Storage, id string) bool {
	_, ok := rooms[id]
	return ok
}
