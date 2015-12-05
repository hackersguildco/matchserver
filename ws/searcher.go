package ws

import "github.com/cheersapp/matchserver/utils"

type searcher struct {
	directory  map[string]*actor
	search     chan *searchActor
	register   chan *registerActor
	unregister chan string
}

// Searcher for the actors on the system
// TODO put this in a better place
var SearcherVar = searcher{
	directory:  make(map[string]*actor),
	search:     make(chan *searchActor, 256),
	register:   make(chan *registerActor, 256),
	unregister: make(chan string, 256),
}

func (s *searcher) Run() {
	for {
		select {
		case search := <-s.search:
			actorRef, _ := s.directory[search.name]
			search.response <- actorRef
			close(search.response)
		case register := <-s.register:
			// creates or find an actor
			actorRef, ok := s.directory[register.name]
			utils.Log.Infof("Looking for actor: %s --- %v", register.name, actorRef)
			if !ok {
				actorRef = createActor(register.name)
				go actorRef.run()
				go actorRef.startTimer()
				s.directory[register.name] = actorRef
			}
			register.response <- actorRef
			close(register.response)
		case username := <-s.unregister:
			if _, ok := s.directory[username]; ok {
				delete(s.directory, username)
			}
		}
	}
}
