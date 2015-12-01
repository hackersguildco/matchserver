package ws

type searcher struct {
	directory  map[string]*actor
	search     chan *searchActor
	register   chan *registerActor
	unregister chan string
}

// Searcher for the actors on the system
// TODO put this in a better place
var searcherVar = searcher{
	directory:  make(map[string]*actor),
	search:     make(chan *searchActor),
	register:   make(chan *registerActor),
	unregister: make(chan string, 256),
}

func (s *searcher) Run() {
	for {
		select {
		case search := <-s.search:
			actorVar, _ := s.directory[search.username]
			search.response <- actorVar
			close(search.response)
		case register := <-s.register:
			// creates or find an actor
			actorVar, ok := s.directory[register.username]
			if !ok {
				actorVar = createActor(register.username)
				go actorVar.run()
				s.directory[register.username] = actorVar
			}
			register.response <- actorVar
			close(register.response)
		case username := <-s.unregister:
			if _, ok := s.directory[username]; ok {
				delete(s.directory, username)
			}
		}
	}
}

func createActor(name string) *actor {
	return &actor{
		name:        name,
		connections: []*connection{},
		systemHub:   &System,
		strokes:     make(chan *postStroke),
		responses:   make(chan []byte),
		nearUsers:   make(chan []string, 256),
		ping:        make(chan *actor),
		pong:        make(chan *actor),
	}
}
