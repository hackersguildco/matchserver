package ws

// Represents an user in the system that is doing a cheers
type actor struct {
	name        string
	connections []*connection
	systemHub   *systemHub
	strokes     chan *postStroke
	nearUsers   chan []string
	responses   chan []byte
	ping        chan *actor
	pong        chan *actor
}

func (a *actor) run() {
	//TODO close resources
	for {
		select {
		case postStrokeVar, more := <-a.strokes:
			if more {
				a.persist(postStrokeVar)
			}
		case users, more := <-a.nearUsers:
			if more {
				//TODO do something with the users ping pong
			}
		}

	}
}

func (a *actor) persist(postStrokeVar *postStroke) {
	persistorVar := persistor{
		persist:  make(chan *postStroke),
		response: a.nearUsers,
	}
	go persistorVar.run()
	persistorVar.persist <- postStrokeVar
}

func (a *actor) pingPong() {
}
