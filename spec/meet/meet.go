package meetspec

import "github.com/gorilla/websocket"

type SingleMeetInfo struct {
	ID      string
	Owner   string
	Clients []*SingleMeetClient
	Host    *SingleMeetClient
}

type SingleMeetClient struct {
	Email           string
	Control         *websocket.Conn
	AudioConn       *websocket.Conn
	VideoConnScreen *websocket.Conn
	//VideoConnCamera websocket.Conn
}
