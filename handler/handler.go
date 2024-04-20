// handler.go

package handler

import (
	"encoding/json"
	"log"
	"meatsrv/utils"
	"net/http"
)

// HandlerConfig represents the configuration for the handler.
type HandlerConfig struct {
	Decoder  func(r *http.Request) (interface{}, error)
	Endpoint func(req interface{}) (interface{}, error)
}

type WebsocketHandlerConfig struct{
	Decoder  func(r *http.Request) (interface{}, error)
	Endpoint func(w http.ResponseWriter,req interface{})
}

// NewHandler creates a new handler with the given configuration.
func NewHandler(config HandlerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := config.Decoder(r)
		if err != nil {
			utils.WriteError(w, err, http.StatusBadRequest)
			return
		}
		resp, err := config.Endpoint(req)
		if err != nil {
			utils.WriteError(w, err, http.StatusInternalServerError)
			return
		}
		b, marshalError := json.Marshal(resp)
		if marshalError != nil {
			log.Println("Error encoding to json: ", marshalError.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
}

func NewWebsocketHandler(config WebsocketHandlerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := config.Decoder(r)
		if err != nil {
			utils.WriteError(w, err, http.StatusBadRequest)
			return
		}
		config.Endpoint(w,req)
	})
}