package meet

import (
	"log"
	"meatsrv/meetops"
	"meatsrv/session"
	meetspec "meatsrv/spec/meet"
	"meatsrv/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

type BL struct {
	sessionManager *session.SessionManager
	meetManager *meetops.MeetManager
}

func NewMeetBL(sessionManager *session.SessionManager,meetManager *meetops.MeetManager) *BL {
	return &BL{
		sessionManager: sessionManager,
		meetManager: meetManager,
	}
}

func (svc *BL) StartMeet(w http.ResponseWriter, req meetspec.StartMeetRequest) {
	//Vaildate otp
	session, err := svc.sessionManager.GetSessionFromOTP(req.OTP)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	svc.sessionManager.PurgeOTP(req.OTP)
	upgrader := websocket.Upgrader{}
	headers := http.Header{}
	meetid := utils.GenerateRandomString(15)
	headers.Add("Meetid", meetid)
	conn, err := upgrader.Upgrade(w, req.Request, headers)
	if err != nil {
		log.Printf("BL:StartMeet: Upgrade failed for email=%s, error=%s", session.Email, err.Error())
		return
	}
	meetInfo:=meetspec.SingleMeetInfo{
		ID:meetid,
		Owner: session.Email,
		Clients: []*meetspec.SingleMeetClient{
			&meetspec.SingleMeetClient{Email: session.Email,Control: conn},
		},
	}
	svc.meetManager.Add(meetInfo)
	go svc.processMeetLoop(meetInfo)
}

func (svc *BL) processMeetLoop(meetInfo meetspec.SingleMeetInfo){
	for{
		host:=meetInfo.Host
		//TODO
	}
}