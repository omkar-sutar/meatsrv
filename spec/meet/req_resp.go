package meetspec

import "net/http"

type StartMeetRequest struct {
	OTP     string `validate:"required"`
	Request *http.Request
}

type StartMeetResponse struct {
}

type JoinMeetRequest struct {
}

type JoinMeetResponse struct {
}
