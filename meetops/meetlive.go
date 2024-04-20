package meetops

import (
	meetspec "meatsrv/spec/meet"
	"sync"
)

type MeetManager struct {
	meets map[string]meetspec.SingleMeetInfo
	mu    *sync.Mutex
}

func NewMeetManager() *MeetManager {
	return &MeetManager{
		meets: map[string]meetspec.SingleMeetInfo{},
		mu:    &sync.Mutex{},
	}
}

func (s *MeetManager) Add(meetinfo meetspec.SingleMeetInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.meets[meetinfo.ID] = meetinfo
}

func (s *MeetManager) Delete(meetinfo meetspec.SingleMeetInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.meets, meetinfo.ID)
}

func (s *MeetManager) Get(id string) meetspec.SingleMeetInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.meets[id]
}
