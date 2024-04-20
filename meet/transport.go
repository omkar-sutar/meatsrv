package meet

import (
	"meatsrv/handler"
	meetspec "meatsrv/spec/meet"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func AddHandlers(router *mux.Router, endpoints MeetEndpoints) {
	router.Path("/meet/start").Handler(handler.NewWebsocketHandler(handler.WebsocketHandlerConfig{
		Decoder: func(r *http.Request) (interface{}, error) {
			return decodeStartMeetRequest(r)
		},
		Endpoint: endpoints.StartMeetEndpoint.Ep,
	}))
}
func decodeStartMeetRequest(r *http.Request) (meetspec.StartMeetRequest, error) {
	var req meetspec.StartMeetRequest
	req.OTP = r.URL.Query().Get("otp")
	req.Request = r
	validate := validator.New()
	err := validate.Struct(req)
	return req, err
}
