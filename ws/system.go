package ws

type systemHub struct {
	// Registered connections.
	connections map[*connection]bool
	// Register requests from the connections.
	register chan *connection
	// Unregister requests from connections.
	unregister chan *connection
}

// System global for handling connection management
var System = systemHub{
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *systemHub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
				close(c.receive)
			}
		}
	}
}
