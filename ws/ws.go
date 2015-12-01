package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// ServeWS handles websocket requests from the peer.
func ServeWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{
		send:    make(chan []byte, 256),
		receive: make(chan []byte, 256),
		ws:      ws,
		name:    username,
	}
	System.register <- c
	go c.writePump()
	go c.run()
	c.readPump()
}
