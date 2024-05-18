package meetspec

import "net/http"

type StartMeetRequest struct {
	OTP     string `validate:"required"`
	Request *http.Request
}

type StartMeetResponse struct {
	MeetID string `json:"meetid"`
}

type JoinMeetRequest struct {
	OTP     string `validate:"required"`
	Request *http.Request
	MeetID  string `json:"meetid" validate:"required"`
}

type JoinMeetResponse struct {
}
