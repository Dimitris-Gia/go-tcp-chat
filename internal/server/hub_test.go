package server

import (
	"net"
	"testing"
	"time"
)

func TestRegisterClient_AddsConnection(t *testing.T) {
	hub := NewHub()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client := NewClient(serverConn, "testuser")
	hub.Register(client)

	if got, want := hub.ClientCount(), 1; got != want {
		t.Fatalf("expected %d registered client, got %d", want, got)
	}
}

func TestUnregisterClient_RemovesConnection(t *testing.T) {
	hub := NewHub()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client := NewClient(serverConn, "testuser")
	hub.Register(client)
	hub.Unregister(client)

	if got, want := hub.ClientCount(), 0; got != want {
		t.Fatalf("expected %d registered clients, got %d", want, got)
	}
}

func TestHub_ClientCount_MultipleClients(t *testing.T) {
	hub := NewHub()
	conns := make([]net.Conn, 3)
	for i := range conns {
		s, c := net.Pipe()
		defer s.Close()
		defer c.Close()
		hub.Register(NewClient(s, "user"))
		conns[i] = s
	}
	if got := hub.ClientCount(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

// readFrom drains one message from the client side of a net.Pipe within 1s.
func readFrom(t *testing.T, conn net.Conn) string {
	t.Helper()
	conn.SetDeadline(time.Now().Add(time.Second))
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	return string(buf[:n])
}

func TestHub_BroadcastAll(t *testing.T) {
	hub := NewHub()

	s1, c1 := net.Pipe()
	s2, c2 := net.Pipe()
	defer s1.Close(); defer c1.Close()
	defer s2.Close(); defer c2.Close()

	cl1 := NewClient(s1, "a")
	cl2 := NewClient(s2, "b")
	go cl1.WritePump()
	go cl2.WritePump()
	hub.Register(cl1)
	hub.Register(cl2)

	hub.BroadcastAll("hello\n")

	if msg := readFrom(t, c1); msg != "hello\n" {
		t.Fatalf("cl1 expected 'hello\\n', got %q", msg)
	}
	if msg := readFrom(t, c2); msg != "hello\n" {
		t.Fatalf("cl2 expected 'hello\\n', got %q", msg)
	}
}

func TestHub_BroadcastExcept_SkipsSender(t *testing.T) {
	hub := NewHub()

	s1, c1 := net.Pipe()
	s2, c2 := net.Pipe()
	defer s1.Close(); defer c1.Close()
	defer s2.Close(); defer c2.Close()

	cl1 := NewClient(s1, "sender")
	cl2 := NewClient(s2, "receiver")
	go cl1.WritePump()
	go cl2.WritePump()
	hub.Register(cl1)
	hub.Register(cl2)

	hub.BroadcastExcept("msg\n", cl1)

	// cl2 should receive the message
	if msg := readFrom(t, c2); !containsStr(msg, "msg") {
		t.Fatalf("cl2 expected message containing 'msg', got %q", msg)
	}

	// cl1 should NOT receive anything within a short window
	c1.SetDeadline(time.Now().Add(100 * time.Millisecond))
	buf := make([]byte, 64)
	if n, _ := c1.Read(buf); n > 0 {
		t.Fatalf("sender should not receive its own broadcast, got %q", string(buf[:n]))
	}
}

func TestHub_RepromptAll(t *testing.T) {
	hub := NewHub()

	s1, c1 := net.Pipe()
	defer s1.Close(); defer c1.Close()

	cl1 := NewClient(s1, "alice")
	go cl1.WritePump()
	hub.Register(cl1)

	hub.RepromptAll()

	msg := readFrom(t, c1)
	if !containsStr(msg, "[alice]:") {
		t.Fatalf("expected prompt containing '[alice]:', got %q", msg)
	}
}

func TestHub_RepromptExcept_SkipsExcluded(t *testing.T) {
	hub := NewHub()

	s1, c1 := net.Pipe()
	s2, c2 := net.Pipe()
	defer s1.Close(); defer c1.Close()
	defer s2.Close(); defer c2.Close()

	cl1 := NewClient(s1, "alice")
	cl2 := NewClient(s2, "bob")
	go cl1.WritePump()
	go cl2.WritePump()
	hub.Register(cl1)
	hub.Register(cl2)

	hub.RepromptExcept(cl1)

	// cl2 should get a prompt
	if msg := readFrom(t, c2); !containsStr(msg, "[bob]:") {
		t.Fatalf("expected prompt for bob, got %q", msg)
	}

	// cl1 should NOT get a prompt
	c1.SetDeadline(time.Now().Add(100 * time.Millisecond))
	buf := make([]byte, 64)
	if n, _ := c1.Read(buf); n > 0 {
		t.Fatalf("excluded client should not receive reprompt, got %q", string(buf[:n]))
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
