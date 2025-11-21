package ws

import "sync"

type Hub struct {
	clients map[string][]*Client
	mx      sync.RWMutex
}

func NewHub() (*Hub, error) {
	return &Hub{
		clients: make(map[string][]*Client),
		mx:      sync.RWMutex{},
	}, nil
}

type Client struct {
	send chan []byte
}

func (c Client) Read() <-chan []byte {
	return c.send
}
