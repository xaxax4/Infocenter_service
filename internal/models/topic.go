package models

import (
	"sync"
)

type Topic struct {
	Name      string
	Clients   map[*Client]struct{}
	Broadcast chan Message
	mu        sync.RWMutex
	done      chan struct{}
}

// since MiddleMan is in another package, encapsulation is preserved

func NewTopic(name string) *Topic {
	return &Topic{
		Name:      name,
		Clients:   make(map[*Client]struct{}),
		Broadcast: make(chan Message, 10),
		done:      make(chan struct{}),
	}
}

func (t *Topic) Lock()                 { t.mu.Lock() }
func (t *Topic) Unlock()               { t.mu.Unlock() }
func (t *Topic) RLock()                { t.mu.RLock() }
func (t *Topic) RUnlock()              { t.mu.RUnlock() }
func (t *Topic) Close()                { close(t.done) }
func (t *Topic) Done() <-chan struct{} { return t.done }
