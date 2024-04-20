package authenticate

import (
	"errors"
	"log"
	"meatsrv/session"
	authenticatespec "meatsrv/spec/authenticate"
)

type BL struct {
	sessionManager *session.SessionManager
}

func NewAuthenticateBL(sessionManager *session.SessionManager) *BL {
	return &BL{
		sessionManager: sessionManager,
	}
}

func (svc *BL) Authenticate(req authenticatespec.AuthenticateRequest) (*authenticatespec.AuthenticateResponse, error) {
	log.Println("BL:Authenticate email=",req.Email)
	if req.Email == "omkar" && req.Password == "omkar" {
		token := svc.sessionManager.GenerateAccessToken(req.Email)
		return &authenticatespec.AuthenticateResponse{
			AccessToken: token,
		}, nil
	}
	return nil, errors.New("invalid creds")
}

func (svc *BL) GenerateOTP(req authenticatespec.OTPRequest) (*authenticatespec.OTPResponse, error) {
	log.Println("BL:GenerateOTP")
	//Get session
	session,err:=svc.sessionManager.GetSessionFromAccessToken(req.AccessToken)
	if err!=nil{
		return nil,err
	}

	otp,err:=svc.sessionManager.GenerateOTP(req.AccessToken)
	if err!=nil{
		return nil,err
	}
	log.Printf("OTP generated for email=%s\n",session.Email)
	return &authenticatespec.OTPResponse{
		OTP: otp,
	}, nil
}
