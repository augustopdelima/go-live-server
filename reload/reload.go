package reload

import "sync"

type Reloader struct {
	mu sync.RWMutex

	clients map[chan struct{}]struct{}
}

func New() *Reloader {
	return &Reloader{
		clients: make(map[chan struct{}]struct{}),
	}
}

func (r *Reloader) Subscribe() chan struct{} {
	client := make(chan struct{},1)

	r.mu.Lock()
	r.clients[client] = struct{}{}
	r.mu.Unlock()

	return  client
}

func (r *Reloader) Unsubscribe(client chan struct{}) {
	r.mu.Lock()
	delete(r.clients,client)
	r.mu.Unlock()

	close(client)
}

func (r *Reloader) Broadcast() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		select {
			case client <- struct{}{}:
			default:
		}
	}
}
