package ws

type searchActor struct {
	name     string
	response chan *actor
}

type registerActor struct {
	name     string
	response chan *actor
}

type postStroke struct {
	userID string
	Info   string
	Loc    []float64
}
