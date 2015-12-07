package test

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"github.com/cheersapp/matchserver/ws"

	. "github.com/onsi/gomega"
)

func createPostStroke(info string, loc []float64) (*ws.PostStroke, []byte) {
	stroke := ws.PostStroke{
		Info: info,
		Loc:  loc,
	}
	json, _ := json.Marshal(stroke)
	return &stroke, json
}

func BeIn(arr []interface{}, val interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func matchOtherTwo(wsConn *websocket.Conn, infoA, infoB string) {
	_, resp12, err12 := wsConn.ReadMessage()
	_, resp13, err13 := wsConn.ReadMessage()
	posibilities := []interface{}{infoA, infoB}
	Expect(err12).To(BeNil())
	Expect(err13).To(BeNil())
	Expect(string(resp12)).NotTo(BeEquivalentTo(string(resp13)))
	Expect(BeIn(posibilities, string(resp12))).To(BeTrue())
	Expect(BeIn(posibilities, string(resp13))).To(BeTrue())
}
