package test

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

func createClient(urlString, username string) *websocket.Conn {
	url, err := url.Parse(urlString)
	if err != nil {
		panic("malformed url")
	}
	dialer := websocket.Dialer{}
	newURL := fmt.Sprintf("ws://%s/ws/%s", url.Host, username)
	conn, _, err := dialer.Dial(newURL, nil)
	if err != nil {
		panic("client couln't be opened")
	}
	return conn
}
