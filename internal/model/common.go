package model

import "sync"

type Status int

const (
	StatusInit Status = iota
	StatusStop
	StatusPlayDone
)

type StatusController struct {
	status Status
	sync.Mutex
}

func (s *StatusController) Load() Status {
	s.Lock()
	defer s.Unlock()
	return s.status
}
func (s *StatusController) SetStop() {
	s.Lock()
	defer s.Unlock()
	switch s.status {
	case StatusInit:
	default:
		return
	}
	s.status = StatusStop
}
func (s *StatusController) SetPlayDone() {
	s.Lock()
	defer s.Unlock()
	switch s.status {
	case StatusInit:
	default:
		return
	}
	s.status = StatusPlayDone
}
