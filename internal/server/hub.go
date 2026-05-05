package server

import (
	"sync"

	"net-cat/internal/logging"
)

// Hub manages all connected chat clients.
type Hub struct {
	mu     sync.Mutex
	clients map[*Client]struct{}
	logger  *logging.Logger
}

// NewHub creates an empty Hub with an optional logger.
func NewHub(logger *logging.Logger) *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		logger:  logger,
	}
}

// Register adds a client to the hub.
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = struct{}{}
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
	if h.logger != nil {
		h.logger.LogEvent(logging.LevelInfo, logging.EventClientDisconnected, logging.ClientDisconnectedData(client.GetName()))
	}
}

// IsNameTaken returns true if any connected client already has the given name.
func (h *Hub) IsNameTaken(name string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		if c.GetName() == name {
			return true
		}
	}
	return false
}

// ClientCount returns the number of currently connected clients.
func (h *Hub) ClientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

// BroadcastExcept delivers msg to every client except the excluded one.
func (h *Hub) BroadcastExcept(msg string, exclude *Client) {
	h.mu.Lock()
	targets := make([]*Client, 0, len(h.clients))
	for c := range h.clients {
		if c != exclude {
			targets = append(targets, c)
		}
	}
	h.mu.Unlock()

	for _, c := range targets {
		c.Deliver("\n" + msg)
	}
}

// BroadcastAll delivers msg to every connected client.
func (h *Hub) BroadcastAll(msg string) {
	h.mu.Lock()
	targets := make([]*Client, 0, len(h.clients))
	for c := range h.clients {
		targets = append(targets, c)
	}
	h.mu.Unlock()

	for _, c := range targets {
		c.Deliver(msg)
	}
}

// RepromptExcept sends each client their current prompt, except the excluded one.
func (h *Hub) RepromptExcept(exclude *Client) {
	h.mu.Lock()
	targets := make([]*Client, 0, len(h.clients))
	for c := range h.clients {
		if c != exclude {
			targets = append(targets, c)
		}
	}
	h.mu.Unlock()

	for _, c := range targets {
		c.Deliver(c.Prompt())
	}
}

// RepromptAll sends every connected client their current prompt.
func (h *Hub) RepromptAll() {
	h.mu.Lock()
	targets := make([]*Client, 0, len(h.clients))
	for c := range h.clients {
		targets = append(targets, c)
	}
	h.mu.Unlock()

	for _, c := range targets {
		c.Deliver(c.Prompt())
	}
}
