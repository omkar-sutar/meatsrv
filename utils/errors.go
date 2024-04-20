package utils

import (
	"encoding/json"
	"log"
	errorspec "meatsrv/spec/error"
	"net/http"
)

func WriteError(w http.ResponseWriter,err error, status int) {
	w.WriteHeader(status)
	errResp := errorspec.ErrorResp{Error: err.Error()}
	b, marshalError := json.Marshal(errResp)
	if marshalError != nil {
		log.Println("Error encoding to json: ", marshalError.Error())
		return
	}
	w.Write(b)
}
