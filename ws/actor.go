package ws

import (
	"time"

	"github.com/cheersapp/matchserver/utils"
)

const timeAlive = 3

// Represents an user in the system that is doing a cheers
type actor struct {
	name             string
	info             []byte
	timer            *time.Timer
	connections      []*connection
	pongedActors     map[*actor]bool
	addConnection    chan *connection
	removeConnection chan *connection
	strokes          chan *PostStroke
	nearUsers        chan []string
	responses        chan []byte
	ping             chan *actor
	pong             chan *actor
	poisonPill       chan bool
}

func createActor(name string) *actor {
	return &actor{
		name:             name,
		connections:      []*connection{},
		addConnection:    make(chan *connection),
		removeConnection: make(chan *connection),
		pongedActors:     make(map[*actor]bool),
		strokes:          make(chan *PostStroke),
		responses:        make(chan []byte),
		nearUsers:        make(chan []string, 256),
		ping:             make(chan *actor, 256),
		pong:             make(chan *actor, 256),
		poisonPill:       make(chan bool, 1),
	}
}

func (a *actor) run() {
	//TODO close resources
	utils.Log.Infof("Running actor: %s", a.name)
	for {
		select {
		case _, more := <-a.poisonPill:
			if more {
				a.die()
			}
		case conn := <-a.addConnection:
			a.connections = append(a.connections, conn)
		case conn := <-a.removeConnection:
			a.removeConnectionBy(conn)
		case postStrokeVar, more := <-a.strokes:
			if more {
				a.persist(postStrokeVar)
			}
		case users, more := <-a.nearUsers:
			if more {
				for _, u := range users {
					searchActorVar := searchActor{
						name:     u,
						response: make(chan *actor),
					}
					SearcherVar.search <- &searchActorVar
					actorRef := <-searchActorVar.response
					actorRef.ping <- a
					a.pongedActors[actorRef] = false
				}
			}
		case actorPing, more := <-a.ping:
			if more {
				if _, ok := a.pongedActors[actorPing]; !ok {
					actorPing.pong <- a
					a.broadcast(actorPing.info)
				}
			}
		case actorPong, more := <-a.pong:
			if more {
				a.pongedActors[actorPong] = true
				actorPong.timer.Stop()
				a.broadcast(actorPong.info)
			}
		}
	}
}

// timeout to kill the actor
func (a *actor) startTimer() {
	a.timer = time.NewTimer(timeAlive * time.Second)
	<-a.timer.C
	a.poisonPill <- true
}

// sends the persist message
func (a *actor) persist(postStrokeVar *PostStroke) {
	persistorVar := persistor{
		persist:  make(chan *PostStroke),
		response: a.nearUsers,
	}
	a.info = []byte(postStrokeVar.Info)
	go persistorVar.run()
	persistorVar.persist <- postStrokeVar
}

// sends data to all connectios
func (a *actor) broadcast(info []byte) {
	for _, conn := range a.connections {
		conn.send <- info
	}
}

// removes a connection
func (a *actor) removeConnectionBy(conn *connection) {
	for i, c := range a.connections {
		if c == conn {
			a.connections = append(a.connections[:i], a.connections[i+1:]...)
			return
		}
	}
}

// finish an actor
func (a *actor) die() {
	// kills the referenced actors
	utils.Log.Infof("Actor dying: %s -- with %d connections", a.name, len(a.connections))
	close(a.responses)
	for actorRef, ponged := range a.pongedActors {
		if ponged {
			actorRef.poisonPill <- true
		}
	}
	// closes all the actor connections
	for _, conn := range a.connections {
		utils.Log.Infof("Closing connection: %s", conn.name)
		close(conn.send)
	}
	SearcherVar.unregister <- a.name
}
