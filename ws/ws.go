package ws

import (
	"log"
	"net/http"

	"github.com/hackersguildco/matchserver/utils"

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
	c := createConnection(username, ws)
	utils.Log.Infof("Creating connection: %s", username)
	go c.writePump()
	// go c.processMessages()
	c.readPump()
}
