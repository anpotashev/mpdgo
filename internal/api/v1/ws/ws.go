package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		return
	}
	defer conn.Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		err = conn.WriteMessage(mt, msg)
		if err != nil {
			break
		}
	}
}
