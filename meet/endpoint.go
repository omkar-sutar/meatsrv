package meet

import (
	meetspec "meatsrv/spec/meet"
	"meatsrv/utils"
	"net/http"
)

type MeetEndpoints struct {
	StartMeetEndpoint utils.WebsocketEndpoint
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
	return MeetEndpoints{
		StartMeetEndpoint: startMeetEndpoint,
	}
}
