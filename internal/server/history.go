package server

import "sync"

// History stores all chat messages in order.
type History struct {
	mu      sync.Mutex
	entries []string
}

// NewHistory creates a new empty History.
func NewHistory() *History {
	return &History{
		entries: make([]string, 0),
	}
}

// Add appends a message to the history in a thread-safe way.
func (h *History) Add(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = append(h.entries, msg)
}

// GetAll returns a copy of all history entries in a thread-safe way.
func (h *History) GetAll() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]string, len(h.entries))
	copy(result, h.entries)
	return result
}
