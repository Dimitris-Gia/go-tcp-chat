package server

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Client represents a connected chat client.
type Client struct {
	Conn net.Conn
	Name string
	send chan string
	mu   sync.Mutex
	done bool
}

// NewClient creates a new Client with a buffered send channel.
func NewClient(conn net.Conn, name string) *Client {
	return &Client{
		Conn: conn,
		Name: name,
		send: make(chan string, 256),
	}
}

// WritePump reads from the send channel and writes to the connection.
// Must be run in its own goroutine.
func (c *Client) WritePump() {
	for msg := range c.send {
		c.Conn.Write([]byte(msg))
	}
}

// Deliver sends a message to the client's send channel safely.
// Drops the message silently if the channel is full or the client is closed.
func (c *Client) Deliver(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.done {
		select {
		case c.send <- msg:
		default:
			// drop message if channel is full
		}
	}
}

// Prompt returns the formatted input prompt for this client.
func (c *Client) Prompt() string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), c.GetName())
}

// GetName returns the client's current name safely.
func (c *Client) GetName() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Name
}

// SetName updates the client's name safely.
func (c *Client) SetName(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Name = name
}

// Close marks the client as done and closes the send channel,
// which causes WritePump to exit.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.done {
		c.done = true
		close(c.send)
	}
}
