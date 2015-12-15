package ws

func init() {
	// InitSearcher run it
	searcherVar = &searcher{
		directory:  make(map[string]*actor),
		search:     make(chan *searchActor, 256),
		register:   make(chan *registerActor, 256),
		unregister: make(chan string, 256),
	}
	go searcherVar.Run()
}
