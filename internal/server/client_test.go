package server

import (
	"net"
	"strings"
	"testing"
	"time"
)

func TestNewClient_Fields(t *testing.T) {
	conn, other := net.Pipe()
	defer conn.Close()
	defer other.Close()

	c := NewClient(conn, "alice")
	if c.Name != "alice" {
		t.Fatalf("expected name 'alice', got %q", c.Name)
	}
	if c.Conn != conn {
		t.Fatal("Conn field not set correctly")
	}
}

func TestClient_DeliverAndWritePump(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	c := NewClient(serverConn, "alice")
	go c.WritePump()

	c.Deliver("hello\n")
	c.Close()

	buf := make([]byte, 64)
	clientConn.SetDeadline(time.Now().Add(time.Second))
	n, _ := clientConn.Read(buf)
	if string(buf[:n]) != "hello\n" {
		t.Fatalf("expected 'hello\\n', got %q", string(buf[:n]))
	}
}

func TestClient_Close_Idempotent(t *testing.T) {
	conn, other := net.Pipe()
	defer conn.Close()
	defer other.Close()

	c := NewClient(conn, "alice")
	go c.WritePump()

	// Calling Close twice must not panic
	c.Close()
	c.Close()
}

func TestClient_DeliverAfterClose_DoesNotPanic(t *testing.T) {
	conn, other := net.Pipe()
	defer conn.Close()
	defer other.Close()

	c := NewClient(conn, "alice")
	go c.WritePump()
	c.Close()

	// Must not panic or block
	c.Deliver("should be dropped\n")
}

func TestClient_Prompt_Format(t *testing.T) {
	conn, other := net.Pipe()
	defer conn.Close()
	defer other.Close()

	c := NewClient(conn, "alice")
	prompt := c.Prompt()

	if !strings.HasSuffix(prompt, "[alice]:") {
		t.Fatalf("prompt should end with '[alice]:', got %q", prompt)
	}
	// Should contain a timestamp-like prefix
	if !strings.HasPrefix(prompt, "[20") {
		t.Fatalf("prompt should start with a year like '[20', got %q", prompt)
	}
}
