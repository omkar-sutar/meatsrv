package session

func (s *SessionManager) GenerateOTP(accessToken string) (string, error) {
	session, err := s.Get(accessToken, false)
	if err != nil {
		return "", err
	}
	//Create OTP session
	return s.Add(session.Email, true), nil
}

func (s *SessionManager) GenerateAccessToken(email string) string {
	return s.Add(email, false)
}

func (s *SessionManager) GetSessionFromAccessToken(accessToken string) (*SessionValue, error) {
	return s.Get(accessToken, false)
}

func (s *SessionManager) GetSessionFromOTP(otp string) (*SessionValue, error) {
	return s.Get(otp, true)
}

func (s *SessionManager) PurgeOTP(otp string) {
	s.Delete(otp, true)
}
