package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, _ := upgrader.Upgrade(w, r, nil)

	for {

		_, msg, err := conn.ReadMessage()

		if err != nil {
			break
		}

		conn.WriteMessage(websocket.TextMessage, msg)

	}

}