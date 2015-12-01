package ws

import "fmt"

type persistor struct {
	persist  chan *postStroke
	response chan []string
}

func (p *persistor) run() {
	defer close(p.persist)
	persist, more := <-p.persist
	if more {
		//TODO save the record in the database and send near users
		fmt.Println(persist)
	} else {
		return
	}
}
