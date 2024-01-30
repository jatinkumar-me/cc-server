package ccserver

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
    http.Handle("/", websocket.Handler())
}
