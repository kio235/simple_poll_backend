package main

import (
	"encoding/json"
	"sync"
)

type Broker struct {
	clients map[chan string]struct{}
	mu      sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[chan string]struct{}),
	}
}

func (b *Broker) AddClient(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[ch] = struct{}{}
}

func (b *Broker) RemoveClient(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, ch)
}

func (b *Broker) Broadcast(data interface{}) {
	jsonData, _ := json.Marshal(data)
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		ch <- string(jsonData)
	}
}
