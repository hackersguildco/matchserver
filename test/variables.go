package test

import (
	"github.com/cheersapp/matchserver/ws"
	"github.com/gorilla/websocket"
)

var (
	username1           = "a1"
	wsConnUser1         *websocket.Conn
	postStrokeUser1Byte []byte
	postStrokeUser1     *ws.PostStroke
	username2           = "a2"
	wsConnUser2         *websocket.Conn
	postStrokeUser2Byte []byte
	postStrokeUser2     *ws.PostStroke
	username3           = "a3"
	wsConnUser3         *websocket.Conn
	postStrokeUser3Byte []byte
	postStrokeUser3     *ws.PostStroke
)
