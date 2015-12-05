package test

import (
	"encoding/json"

	"github.com/cheersapp/matchserver/ws"
)

func createPostStroke(info string, loc []float64) (*ws.PostStroke, []byte) {
	stroke := ws.PostStroke{
		Info: info,
		Loc:  loc,
	}
	json, _ := json.Marshal(stroke)
	return &stroke, json
}
