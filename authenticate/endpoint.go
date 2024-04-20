package authenticate

import (
	"errors"
	authenticatespec "meatsrv/spec/authenticate"
	"meatsrv/utils"
)

type AuthenticateEndpoints struct {
	Authenticate utils.Endpoint
	OTP          utils.Endpoint
}

func AddEndpoints(bl *BL) AuthenticateEndpoints {
	authenticateEndpoint := utils.Endpoint{
		Ep: func(req interface{}) (interface{}, error) {
			// Type assertion to authenticate request
			authReq, ok := req.(authenticatespec.AuthenticateRequest)
			if !ok {
				return nil, errors.New("invalid request type")
			}
			return bl.Authenticate(authReq)
		},
	}
	otpEndpoint := utils.Endpoint{
		Ep: func(req interface{}) (interface{}, error) {
			otpReq, ok := req.(authenticatespec.OTPRequest)
			if !ok {
				return nil, errors.New("invalid request type")
			}
			return bl.GenerateOTP(otpReq)
		},
	}
	endpoints := AuthenticateEndpoints{
		Authenticate: authenticateEndpoint,
		OTP:          otpEndpoint,
	}
	return endpoints
}
