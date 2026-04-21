package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"net-cat/internal/connectionhandling"
	"net-cat/internal/parser"
	"net-cat/internal/server"
)

// startServer spins up a real TCP listener on a random port and returns the address.
func startServer(t *testing.T) string {
	t.Helper()

	hub := server.NewHub()
	history := server.NewHistory()

	listener, err := net.Listen("tcp", ":0") // :0 = OS picks a free port
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	t.Cleanup(func() { listener.Close() })

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return // listener closed
			}
			go connectionhandling.HandleConnection(conn, hub, history)
		}
	}()

	return listener.Addr().String()
}

// TestMain_DefaultPort verifies GetPortNumber returns 8989 when no port is given.
func TestMain_DefaultPort(t *testing.T) {
	port, err := parser.GetPortNumber([]string{"./TCPChat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if port != 8989 {
		t.Fatalf("expected default port 8989, got %d", port)
	}
}

// TestMain_InvalidArgs verifies GetPortNumber errors on too many arguments.
func TestMain_InvalidArgs(t *testing.T) {
	_, err := parser.GetPortNumber([]string{"./TCPChat", "2525", "localhost"})
	if err == nil {
		t.Fatal("expected error for too many args, got nil")
	}
	if !strings.Contains(err.Error(), "[USAGE]: ./TCPChat $port") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestMain_ServerAcceptsConnection verifies the server accepts a TCP connection
// and sends the welcome message.
func TestMain_ServerAcceptsConnection(t *testing.T) {
	addr := startServer(t)

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(2 * time.Second))
	r := bufio.NewReader(conn)

	var buf strings.Builder
	for {
		part, err := r.ReadString(':')
		buf.WriteString(part)
		if strings.Contains(buf.String(), "[ENTER YOUR NAME]:") {
			break
		}
		if err != nil {
			t.Fatalf("error reading welcome: %v\ngot so far: %q", err, buf.String())
		}
	}

	if !strings.Contains(buf.String(), "Welcome to TCP-Chat!") {
		t.Fatalf("expected welcome message, got %q", buf.String())
	}
}

// TestMain_ServerHandlesMultipleConnections verifies the server can handle
// multiple simultaneous clients.
func TestMain_ServerHandlesMultipleConnections(t *testing.T) {
	addr := startServer(t)

	connect := func(name string) net.Conn {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			t.Fatalf("could not connect: %v", err)
		}
		r := bufio.NewReader(conn)
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		// drain until "[ENTER YOUR NAME]:" (no trailing newline — read until ':')
		var wb strings.Builder
		for {
			part, err := r.ReadString(':')
			wb.WriteString(part)
			if strings.Contains(wb.String(), "[ENTER YOUR NAME]:") {
				break
			}
			if err != nil {
				break
			}
		}
		fmt.Fprintf(conn, "%s\n", name)
		return conn
	}

	c1 := connect("alice")
	defer c1.Close()
	c2 := connect("bob")
	defer c2.Close()

	// Both connected — server should still be running (no panic/crash)
	// Verify by connecting a third client
	c3 := connect("carol")
	defer c3.Close()
}

// TestMain_ServerBroadcastsMessages verifies a message from one client
// reaches another through the real TCP server.
func TestMain_ServerBroadcastsMessages(t *testing.T) {
	addr := startServer(t)

	// helper: connect, drain welcome, send name, drain until own prompt appears.
	// The server prompt ends with "[name]:" and has NO trailing newline,
	// so we read byte-by-byte and stop when the suffix matches.
	join := func(name string) (net.Conn, *bufio.Reader) {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			t.Fatalf("dial failed: %v", err)
		}
		r := bufio.NewReader(conn)
		conn.SetDeadline(time.Now().Add(3 * time.Second))

		// drain until "[ENTER YOUR NAME]:" — no trailing newline, so read until ':'
		var welcome strings.Builder
		for {
			part, err := r.ReadString(':')
			welcome.WriteString(part)
			if strings.Contains(welcome.String(), "[ENTER YOUR NAME]:") {
				break
			}
			if err != nil {
				break
			}
		}
		fmt.Fprintf(conn, "%s\n", name)

		// drain until own prompt "[name]:" (no trailing newline)
		suffix := fmt.Sprintf("[%s]:", name)
		var buf strings.Builder
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}
			buf.WriteByte(b)
			if strings.HasSuffix(buf.String(), suffix) {
				break
			}
		}
		conn.SetDeadline(time.Time{}) // clear deadline for subsequent reads
		return conn, r
	}

	c1, _ := join("alice")
	defer c1.Close()
	c2, r2 := join("bob")
	defer c2.Close()

	// alice sends a message
	fmt.Fprint(c1, "hello from alice\n")

	// bob should receive it
	c2.SetDeadline(time.Now().Add(3 * time.Second))
	var got strings.Builder
	for {
		line, err := r2.ReadString('\n')
		got.WriteString(line)
		if strings.Contains(got.String(), "hello from alice") {
			break
		}
		if err != nil {
			t.Fatalf("bob did not receive alice's message: %v\ngot: %q", err, got.String())
		}
	}
}
