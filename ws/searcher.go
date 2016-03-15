package ws

import "github.com/hackersguildco/matchserver/utils"

type searcher struct {
	directory  map[string]*actor
	search     chan *searchActor
	register   chan *registerActor
	unregister chan string
}

// Searcher for the actors on the system
var searcherVar *searcher

// Run the searcher
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
				s.directory[register.name] = actorRef
				go actorRef.run()
				go actorRef.startTimer()
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
