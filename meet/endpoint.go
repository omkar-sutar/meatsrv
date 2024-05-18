package meet

import (
	meetspec "meatsrv/spec/meet"
	"meatsrv/utils"
	"net/http"
)

type MeetEndpoints struct {
	StartMeetEndpoint utils.WebsocketEndpoint
	JoinMeetEndpoint  utils.WebsocketEndpoint
}

func AddEndpoints(bl *BL) MeetEndpoints {
	startMeetEndpoint := utils.WebsocketEndpoint{
		Ep: func(w http.ResponseWriter, req interface{}) {
			// Type assertion to authenticate request
			startMeetRequest, ok := req.(meetspec.StartMeetRequest)
			if !ok {
				return
			}
			bl.StartMeet(w, startMeetRequest)
		},
	}
	joinMeetEndpoint := utils.WebsocketEndpoint{
		Ep: func(w http.ResponseWriter, req interface{}) {
			// Type assertion to authenticate request
			joinMeetRequest, ok := req.(meetspec.JoinMeetRequest)
			if !ok {
				return
			}
			bl.JoinMeet(w, joinMeetRequest)
		},
	}
	return MeetEndpoints{
		StartMeetEndpoint: startMeetEndpoint,
		JoinMeetEndpoint:  joinMeetEndpoint,
	}
}
