package main
import (
	"io"
	"encoding/json"
	"math/rand"
	"strconv"
)
var CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

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

func MakeId() string {
	id := ""
	for i := 0; i < 8; i++ {
		// All letters are valid. NO need to check
		id += strconv.Itoa(int(CHARS[rand.Intn(len(CHARS))]))
	}
	if _, ok := RoomStorage[id]; ok {
		// This string has already been created
		return MakeId()
	}
	return id
}

func MakePin() string {
	pin := ""
	nums := "1234567890"
	for i := 0; i < 5; i++ {
		pin += strconv.Itoa(int(nums[rand.Intn(len(CHARS))]))
	}
	return pin
}