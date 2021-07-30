package main

import (
	"sync"
)

type ForwardManager struct {
	lns map[string]*Forward

	sync.Mutex
}

func (s *ForwardManager) AddForward(key string, forward *Forward) {
	s.Lock()
	defer s.Unlock()
	s.lns[key] = forward

}

func (s *ForwardManager) RemoveForward(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.lns, key)

}

func (s *ForwardManager) GetForward(key string) *Forward {
	s.Lock()
	defer s.Unlock()
	return s.lns[key]
}
