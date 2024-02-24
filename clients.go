package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"golang.org/x/net/websocket"
)

type Client struct {
	user User
	conn *websocket.Conn
}

func GenerateUser() (user User) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("Error generating random bytes:", err)
	}

	randomId := hex.EncodeToString(b)
	user.UserId = UserId(randomId[:8])
	user.UserName = "Anon" + string(user.UserId)
	return user
}
