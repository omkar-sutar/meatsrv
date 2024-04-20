package utils

import "net/http"

type Endpoint struct {
	Ep func(interface{}) (interface{}, error)
}

type WebsocketEndpoint struct {
	Ep func(http.ResponseWriter,interface{})
}
