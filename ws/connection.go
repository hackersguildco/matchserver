package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cheersappio/matchserver/utils"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn
	// name of the connection
	name string
	// actor reference
	actorRef *actor
	// Buffered channel of outbound messages.
	send chan []byte
	// channel for inbound messages
	receive chan []byte
	// channel to kill the connection
	poisonPill chan bool
}

func createConnection(name string, ws *websocket.Conn) *connection {
	return &connection{
		ws:         ws,
		name:       name,
		send:       make(chan []byte, 256),
		receive:    make(chan []byte, 256),
		poisonPill: make(chan bool),
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		if c.actorRef != nil {
			c.actorRef.removeConnection <- c
		}
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.receive <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		utils.Log.Infof("Finishing connection: %s", c.name)
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-c.poisonPill:
			if len(c.send) == 0 {
				utils.Log.Infof("Closing websocket: %s", c.name)
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			go c.dieLater()
		case message := <-c.send: // channel used to finish the connection when it's closed
			utils.Log.Infof("Into send channel for: %s, message: %s", c.name, string(message))
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case message := <-c.receive:
			c.processMessages(message)
		}
	}
}

// process each message
func (c *connection) processMessages(message []byte) {
	register := registerActor{
		name:     c.name,
		response: make(chan *actor),
	}
	utils.Log.Infof("Creating actor: %s", register.name)
	searcherVar.register <- &register
	actorRef := <-register.response
	actorRef.addConnection <- c
	//creates the postStroke
	postStrokeVar := postStroke{}
	json.Unmarshal(message, &postStrokeVar)
	postStrokeVar.userID = actorRef.name
	actorRef.strokes <- &postStrokeVar
	for {
		// expects all the responses from the actor until it dies
		response, more := <-actorRef.responses
		if more {
			c.send <- response
		} else {
			break
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) dieLater() {
	timer := time.NewTimer(100 * time.Millisecond)
	<-timer.C
	c.poisonPill <- true
}
