package session

import (
	"errors"
	"log"
	"meatsrv/constants"
	"meatsrv/utils"
	"sync"
	"time"
)

type SessionManager struct {
	sessions map[SessionKey]SessionValue
	mu       *sync.Mutex
}

type SessionKey struct {
	KeyType int
	Key     string
}

type SessionValue struct {
	Email    string
	createTS time.Time
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: map[SessionKey]SessionValue{},
		mu:       &sync.Mutex{},
	}
}

func (s *SessionManager) Add(email string, otp bool) (token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	token = utils.GenerateSecureToken()
	sessionKey := SessionKey{
		KeyType: constants.AccessTokenType,
		Key:     token,
	}
	if otp {
		sessionKey.KeyType = constants.OTPTokenType
	}
	s.sessions[sessionKey] = SessionValue{
		Email:    email,
		createTS: time.Now(),
	}
	log.Println("Successfully created session for email=", email)
	return
}

func (s *SessionManager) Delete(token string, otp bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sessionKey := SessionKey{
		KeyType: constants.AccessTokenType,
		Key:     token,
	}
	if otp {
		sessionKey.KeyType = constants.OTPTokenType
	}
	session, ok := s.sessions[sessionKey]
	if !ok {
		return
	}
	delete(s.sessions, sessionKey)
	log.Println("Session ended for email=", session.Email)
}

func (s *SessionManager) Get(token string, otp bool) (*SessionValue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sessionKey := SessionKey{
		KeyType: constants.AccessTokenType,
		Key:     token,
	}
	if otp {
		sessionKey.KeyType = constants.OTPTokenType
	}
	session, ok := s.sessions[sessionKey]

	//OTP type
	if otp {
		if !ok {
			return nil, errors.New(constants.InvalidOTP)
		}
		delete(s.sessions, sessionKey)
		return &session, nil
	}

	//AccessToken type
	if !ok || session.createTS.Add(constants.TokenValidityDuration).Before(time.Now()) {
		delete(s.sessions, sessionKey)
		return nil, errors.New(constants.SessionExpired)
	}
	return &session, nil
}
