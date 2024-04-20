package authenticatespec

type AuthenticateRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

type AuthenticateResponse struct {
	AccessToken string `json:"accessToken"`
}

type OTPRequest struct {
	AccessToken string `validate:"required"`
}

type OTPResponse struct{
	OTP string `json:"otp"`
}
