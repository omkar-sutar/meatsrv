package meetspec

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MaxMeetTime = 30 * time.Minute
)

const (
	BroadcastMessageType = "broadcast"
	DirectMessageType    = "direct"
	ErrorMessageType     = "error"
)

type SingleMeetInfo struct {
	mu           *sync.Mutex
	ID           string
	Clients      []SingleMeetClient
	MessageQueue chan Message
	Done         chan int
}

func NewMeet(meetID string) *SingleMeetInfo {
	return &SingleMeetInfo{
		mu:           &sync.Mutex{},
		ID:           meetID,
		Done:         make(chan int),
		MessageQueue: make(chan Message, 128),
	}
}

func (meet *SingleMeetInfo) DeleteClient(email string) {
	meet.mu.Lock()
	defer meet.mu.Unlock()
	newClients := []SingleMeetClient{}
	for _, client := range meet.Clients {
		if client.Email != email {
			newClients = append(newClients, client)
		}
	}
	meet.Clients = newClients
}

func (meet *SingleMeetInfo) AddClient(client SingleMeetClient) {
	meet.mu.Lock()
	defer meet.mu.Unlock()
	meet.Clients = append(meet.Clients, client)
}

type SingleMeetClient struct {
	Email     string
	Conn      *websocket.Conn
	WriteChan chan Message
	Done      chan int
}

type Message struct {
	To     string
	From   string
	Type   string //E.g. Direct
	Intent string //e.g. 'RTCPeerRequest'
	Body   string
}
