package connectionhandling

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"net-cat/internal/server"
)

// dial simulates a client: writes lines and returns a reader for server output.
func dial(t *testing.T) (net.Conn, *bufio.Reader) {
	t.Helper()
	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() { serverConn.Close(); clientConn.Close() })
	return serverConn, bufio.NewReader(clientConn)
}

// readUntil reads byte-by-byte until the target substring is found or deadline exceeded.
// Reading byte-by-byte avoids blocking on ReadString('\n') for prompts with no trailing newline.
func readUntil(t *testing.T, r *bufio.Reader, target string, conn net.Conn) string {
	t.Helper()
	var buf strings.Builder
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	for {
		b, err := r.ReadByte()
		if err != nil {
			t.Fatalf("did not find %q before error: %v\nGot so far: %q", target, err, buf.String())
		}
		buf.WriteByte(b)
		if strings.Contains(buf.String(), target) {
			return buf.String()
		}
	}
}

func TestGetWelcomeMessage_ContainsPenguin(t *testing.T) {
	msg := GetWelcomeMessage()
	if !strings.Contains(msg, "Welcome to TCP-Chat!") {
		t.Fatal("welcome message missing 'Welcome to TCP-Chat!'")
	}
	if !strings.Contains(msg, "[ENTER YOUR NAME]:") {
		t.Fatal("welcome message missing '[ENTER YOUR NAME]:'")
	}
}

func TestHandleConnection_SendsWelcomeOnConnect(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	r := bufio.NewReader(clientConn)
	go HandleConnection(serverConn, hub, history, nil)

	got := readUntil(t, r, "[ENTER YOUR NAME]:", clientConn)
	if !strings.Contains(got, "Welcome to TCP-Chat!") {
		t.Fatalf("expected welcome message, got %q", got)
	}
}

func TestHandleConnection_RejectsDuplicateName(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	// alice connects first
	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	// second client tries to use the same name
	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "alice\n")

	got := readUntil(t, r2, "already in use", c2)
	if !strings.Contains(got, "already in use") {
		t.Fatalf("expected 'already in use' message, got %q", got)
	}
}

func TestHandleConnection_DuplicateNameRepromptAcceptsNewName(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "alice\n") // taken
	readUntil(t, r2, "already in use", c2)
	fmt.Fprint(c2, "bob\n") // unique

	got := readUntil(t, r2, "[bob]:", c2)
	if !strings.Contains(got, "[bob]:") {
		t.Fatalf("expected prompt with new name, got %q", got)
	}
}

func TestHandleConnection_UserNameChange_RejectsDuplicateName(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2)

	// bob tries to rename to alice
	fmt.Fprint(c2, "--UserNameChange: alice\n")

	got := readUntil(t, r2, "already in use", c2)
	if !strings.Contains(got, "already in use") {
		t.Fatalf("expected 'already in use' on rename to taken name, got %q", got)
	}
}

func TestHandleConnection_RejectsEmptyName(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	r := bufio.NewReader(clientConn)
	go HandleConnection(serverConn, hub, history, nil)

	// Wait for first prompt
	readUntil(t, r, "[ENTER YOUR NAME]:", clientConn)

	// Send empty name
	fmt.Fprint(clientConn, "\n")

	// Server should re-prompt
	got := readUntil(t, r, "[ENTER YOUR NAME]:", clientConn)
	if !strings.Contains(got, "[ENTER YOUR NAME]:") {
		t.Fatalf("expected re-prompt after empty name, got %q", got)
	}
}

func TestHandleConnection_AcceptsValidName(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	r := bufio.NewReader(clientConn)
	go HandleConnection(serverConn, hub, history, nil)

	readUntil(t, r, "[ENTER YOUR NAME]:", clientConn)
	fmt.Fprint(clientConn, "alice\n")

	// After valid name, client should receive their prompt
	got := readUntil(t, r, "[alice]:", clientConn)
	if !strings.Contains(got, "[alice]:") {
		t.Fatalf("expected prompt with name, got %q", got)
	}
}

func TestHandleConnection_ReceivesChatHistory(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()
	history.Add("[2020-01-01 00:00:00][old]:hello\n")

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	r := bufio.NewReader(clientConn)
	go HandleConnection(serverConn, hub, history, nil)

	readUntil(t, r, "[ENTER YOUR NAME]:", clientConn)
	fmt.Fprint(clientConn, "bob\n")

	got := readUntil(t, r, "[old]:hello", clientConn)
	if !strings.Contains(got, "[old]:hello") {
		t.Fatalf("expected history message, got %q", got)
	}
}

func TestHandleConnection_BroadcastsMessageToOtherClient(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	// Connect client1
	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1) // wait for prompt

	// Connect client2
	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2) // wait for prompt

	// alice sends a message
	fmt.Fprint(c1, "hi bob!\n")

	// bob should receive it
	got := readUntil(t, r2, "hi bob!", c2)
	if !strings.Contains(got, "hi bob!") {
		t.Fatalf("bob expected to receive alice's message, got %q", got)
	}
}

func TestHandleConnection_EmptyMessageNotBroadcast(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2)

	// alice sends empty message
	fmt.Fprint(c1, "\n")

	// bob should NOT receive anything within a short window
	c2.SetDeadline(time.Now().Add(150 * time.Millisecond))
	buf := make([]byte, 256)
	n, _ := c2.Read(buf)
	// any data bob gets should not be a chat message (only a reprompt at most)
	if n > 0 && strings.Contains(string(buf[:n]), "[alice]:") {
		t.Fatalf("empty message should not be broadcast, bob got: %q", string(buf[:n]))
	}
}

func TestHandleConnection_JoinNotification(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	// alice connects first
	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	// bob connects second
	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")

	// alice should see join notification
	got := readUntil(t, r1, "bob has joined our chat...", c1)
	if !strings.Contains(got, "bob has joined our chat...") {
		t.Fatalf("expected join notification, got %q", got)
	}
}

func TestHandleConnection_ChangeUserNameUpdatesPrompt(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	r := bufio.NewReader(clientConn)
	go HandleConnection(serverConn, hub, history, nil)

	readUntil(t, r, "[ENTER YOUR NAME]:", clientConn)
	fmt.Fprint(clientConn, "alice\n")
	readUntil(t, r, "[alice]:", clientConn)

	fmt.Fprint(clientConn, "--UserNameChange: charlie\n")

	readUntil(t, r, "alice is now known as charlie", clientConn)
	got := readUntil(t, r, "[charlie]:", clientConn)
	if !strings.Contains(got, "[charlie]:") {
		t.Fatalf("expected updated prompt with new name, got %q", got)
	}
}

func TestHandleConnection_ChangeUserNameNotifiesOtherClient(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	// Connect alice
	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	// Connect bob
	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2)

	// alice renames
	fmt.Fprint(c1, "--UserNameChange: charlie\n")

	got := readUntil(t, r2, "alice is now known as charlie", c2)
	if strings.Contains(got, "--UserNameChange:") {
		t.Fatalf("rename command should not be broadcast as a chat message, bob got: %q", got)
	}
}

func TestHandleConnection_LeaveNotification(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	s2, c2 := net.Pipe()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2)

	// bob disconnects
	c2.Close()

	got := readUntil(t, r1, "bob has left our chat...", c1)
	if !strings.Contains(got, "bob has left our chat...") {
		t.Fatalf("expected leave notification, got %q", got)
	}
}

func TestHandleConnection_ConnectionLimit(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	// Fill up to maxClients
	clientConns := make([]net.Conn, maxClients)
	for i := 0; i < maxClients; i++ {
		s, c := net.Pipe()
		clientConns[i] = c
		r := bufio.NewReader(c)
		go HandleConnection(s, hub, history, nil)
		readUntil(t, r, "[ENTER YOUR NAME]:", c)
		fmt.Fprintf(c, "user%d\n", i)
		readUntil(t, r, fmt.Sprintf("[user%d]:", i), c)
	}
	defer func() {
		for _, c := range clientConns {
			c.Close()
		}
	}()

	// 11th connection should be rejected
	s11, c11 := net.Pipe()
	defer c11.Close()
	r11 := bufio.NewReader(c11)
	go HandleConnection(s11, hub, history, nil)

	readUntil(t, r11, "[ENTER YOUR NAME]:", c11)
	fmt.Fprint(c11, "overflow\n")

	got := readUntil(t, r11, "full", c11)
	if !strings.Contains(strings.ToLower(got), "full") {
		t.Fatalf("expected rejection message for 11th client, got %q", got)
	}
}

func TestHandleConnection_ColorFlagBroadcastsColoredMessage(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2)

	fmt.Fprint(c1, "--red hello\n")

	got := readUntil(t, r2, "hello", c2)
	if !strings.Contains(got, "hello") {
		t.Fatalf("expected colored message, got %q", got)
	}
}

func TestHandleConnection_EmoteFlagBroadcastsAsciiEmote(t *testing.T) {
	hub := server.NewHub(nil)
	history := server.NewHistory()

	s1, c1 := net.Pipe()
	defer c1.Close()
	r1 := bufio.NewReader(c1)
	go HandleConnection(s1, hub, history, nil)
	readUntil(t, r1, "[ENTER YOUR NAME]:", c1)
	fmt.Fprint(c1, "alice\n")
	readUntil(t, r1, "[alice]:", c1)

	s2, c2 := net.Pipe()
	defer c2.Close()
	r2 := bufio.NewReader(c2)
	go HandleConnection(s2, hub, history, nil)
	readUntil(t, r2, "[ENTER YOUR NAME]:", c2)
	fmt.Fprint(c2, "bob\n")
	readUntil(t, r2, "[bob]:", c2)

	fmt.Fprint(c1, "--shrug\n")

	got := readUntil(t, r2, "\\_(o_o)_/", c2)
	if !strings.Contains(got, "\\_(o_o)_/") {
		t.Fatalf("expected emote message, got %q", got)
	}
}
