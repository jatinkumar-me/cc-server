package main

import (
	"net/http"
	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()
	http.Handle("/", websocket.Handler(server.handleWS))
	http.ListenAndServe(":8080", nil)
}
