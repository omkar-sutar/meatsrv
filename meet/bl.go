package meet

import (
	"context"
	"log"
	"meatsrv/meetops"
	"meatsrv/session"
	meetspec "meatsrv/spec/meet"
	"meatsrv/utils"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type BL struct {
	sessionManager *session.SessionManager
	meetManager    *meetops.MeetManager
}

func NewMeetBL(sessionManager *session.SessionManager, meetManager *meetops.MeetManager) *BL {
	return &BL{
		sessionManager: sessionManager,
		meetManager:    meetManager,
	}
}

func (svc *BL) StartMeet(w http.ResponseWriter, req meetspec.StartMeetRequest) {
	//Vaildate otp
	session, err := svc.sessionManager.GetSessionFromOTP(req.OTP)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	upgrader := websocket.Upgrader{}
	headers := http.Header{}
	meetid := utils.GenerateRandomString(15)
	headers.Add("Meetid", meetid)
	conn, err := upgrader.Upgrade(w, req.Request, headers)
	if err != nil {
		log.Printf("BL:StartMeet: Upgrade failed for email=%s, error=%s\n", session.Email, err.Error())
		return
	}
	err = conn.WriteJSON(meetspec.StartMeetResponse{
		MeetID: meetid,
	})
	if err != nil {
		log.Println("BL:StartMeet: Response write failed")
	}
	meetInfo := meetspec.NewMeet(meetid)
	meetInfo.AddClient(meetspec.SingleMeetClient{
		Email:     session.Email,
		Conn:      conn,
		Done:      make(chan int, 5),
		WriteChan: make(chan meetspec.Message, 128),
	})
	svc.meetManager.Add(meetInfo)
	// Start processing first client
	svc.ProcessMeet(meetInfo)
}

func (svc *BL) ProcessMeet(meet *meetspec.SingleMeetInfo) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), meetspec.MaxMeetTime)

	go svc.ProcessClient(meet, meet.Clients[0])
	log.Println("Starting goroutine #2")
	svc.MeetMessageQueueProcessorAsync(ctx, meet)

	SignalClientDone := func() {
		log.Printf("Cleaning up meet=%v", meet)
		for _, client := range meet.Clients {
			client.Done <- 0
		}
		time.Sleep(500 * time.Millisecond)
		close(meet.MessageQueue) //Sleep is necessary to ensure ReadConnectionAsync sends data if any before closing MessageQueue
	}
	// Cleanup when meet is done or timer expires
	//Close all client connections
	select {
	case <-meet.Done:
	case <-ctx.Done():
	}
	SignalClientDone() // Terminates ProcessClient
	cancelFunc()       //Terminates MessageQueueProcessorAsync, Read/WriteConnectionAsync if running
	svc.meetManager.Delete(meet.ID)
}

func (svc *BL) ProcessClient(meet *meetspec.SingleMeetInfo, client meetspec.SingleMeetClient) {

	ctx, cancelFunc := context.WithCancel(context.Background())
	log.Println("Starting goroutine #1")
	defer log.Println("End goroutine #1")
	defer client.Conn.Close()
	log.Println("Starting goroutine #3")
	svc.ReadConnectionAsync(ctx, meet, client)
	log.Println("Starting goroutine #4")
	svc.WriteConnectionAsync(ctx, meet, client)

	<-client.Done
	cancelFunc()
	meet.DeleteClient(client.Email)
}

func (svc *BL) ReadConnectionAsync(ctx context.Context, meet *meetspec.SingleMeetInfo, client meetspec.SingleMeetClient) {
	//Reads messages from client and sends them into Message Queue
	go func() {
		defer log.Println("End goroutine #3")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg := meetspec.Message{}
				err := client.Conn.ReadJSON(&msg)
				if err != nil {
					log.Printf("ReadConnectionAsync err for client=%v, err=%v\n", client, err)
					//Signal client deletion
					client.Done <- 1
					return
				}
				msg.From = client.Email
				//Put the message in message queue
				meet.MessageQueue <- msg
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
}

func (svc *BL) MeetMessageQueueProcessorAsync(ctx context.Context, meet *meetspec.SingleMeetInfo) {
	go func() {
		defer log.Println("End goroutine #2")
		select {
		case msg := <-meet.MessageQueue:
			log.Printf("Message retrieved = %v\n", msg)
			switch msg.Type {
			case meetspec.BroadcastMessageType:
				log.Printf("Broadcast message\n")
				for _, client := range meet.Clients {
					if client.Email != msg.From {
						client.WriteChan <- msg
					}
				}
			case meetspec.DirectMessageType:
				for _, client := range meet.Clients {
					if client.Email == msg.To {
						client.WriteChan <- msg
						break
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}()
}

func (svc *BL) WriteConnectionAsync(ctx context.Context, meet *meetspec.SingleMeetInfo, client meetspec.SingleMeetClient) {
	//Extracts messages from Write queue and writes it to the websocket
	go func() {
		defer log.Println("End goroutine #4")
		for {
			select {
			case msg := <-client.WriteChan:
				// Set write deadline
				//client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("WriteConnectionAsync err for client=%s,err=%s\n", client.Email, err.Error())
					//Signal client deletion
					client.Done <- 1
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (svc *BL) JoinMeet(w http.ResponseWriter, req meetspec.JoinMeetRequest) {
	session, err := svc.sessionManager.GetSessionFromOTP(req.OTP)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	meetid := req.MeetID
	meetInfo, ok := svc.meetManager.Get(meetid)
	if !ok {
		log.Printf("BL:JoinMeet: Failed for email=%s, invalid meetid=%s\n", session.Email, meetid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{}
	headers := http.Header{}
	headers.Add("Meetid", meetid)

	conn, err := upgrader.Upgrade(w, req.Request, headers)
	if err != nil {
		log.Printf("BL:JoinMeet: Upgrade failed for email=%s, error=%s\n", session.Email, err.Error())
		return
	}
	client := meetspec.SingleMeetClient{
		Email:     session.Email,
		Conn:      conn,
		Done:      make(chan int, 5),
		WriteChan: make(chan meetspec.Message),
	}
	meetInfo.AddClient(client)
	svc.ProcessClient(meetInfo, client)
}
