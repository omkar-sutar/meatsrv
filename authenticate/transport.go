package authenticate

import (
	"encoding/json"
	"io"
	authenticatespec "meatsrv/spec/authenticate"
	"net/http"
	"strings"

	handler "meatsrv/handler"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func AddHandlers(router *mux.Router, endpoints AuthenticateEndpoints) {
	// Configure the handler

	// Create a new handler and attach it to the router
	router.Path("/authenticate").Methods("POST").Handler(handler.NewHandler(handler.HandlerConfig{
		Decoder: func(r *http.Request) (interface{}, error) {
			return decodeAuthenticateRequest(r)
		},
		Endpoint: endpoints.Authenticate.Ep,
	},
	))

	router.Path("/otp").Methods("GET").Handler(handler.NewHandler(handler.HandlerConfig{
		Decoder: func(r *http.Request) (interface{}, error) {
			return decodeOTPRequest(r)
		},
		Endpoint: endpoints.OTP.Ep,
	},
	))
}

func decodeOTPRequest(r *http.Request) (authenticatespec.OTPRequest, error) {
	var req authenticatespec.OTPRequest

	reqToken := r.Header.Get("Authorization")
	_, accessToken, _ := strings.Cut(reqToken, "Bearer ")
	req.AccessToken = accessToken
	validate := validator.New()
	err := validate.Struct(req)
	return req, err
}

func decodeAuthenticateRequest(r *http.Request) (authenticatespec.AuthenticateRequest, error) {
	var req authenticatespec.AuthenticateRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return req, err
	}
	r.Body.Close()
	err = json.Unmarshal(b, &req)
	if err != nil {
		return req, err
	}
	validate := validator.New()
	err = validate.Struct(req)
	return req, err
}
