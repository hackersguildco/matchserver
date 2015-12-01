package ws

type searchActor struct {
	username string
	response chan *actor
}

type registerActor struct {
	username string
	response chan *actor
}

type postStroke struct {
	body string
	loc  []float64
}
