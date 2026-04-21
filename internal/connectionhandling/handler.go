package connectionhandling

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"net-cat/internal/server"
)

const maxClients = 10

// HandleConnection manages the full lifecycle of a single client connection.
// It must be run in its own goroutine.
func HandleConnection(conn net.Conn, hub *server.Hub, history *server.History) {
	defer conn.Close()

	// 1. Send welcome message (ASCII art + [ENTER YOUR NAME]: prompt)
	conn.Write([]byte(GetWelcomeMessage()))

	reader := bufio.NewReader(conn)

	// 2. Name registration loop — keep asking until a non-empty name is given
	var name string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return // client disconnected before entering a name
		}
		name = strings.TrimSpace(line)
		if name != "" {
			break
		}
		conn.Write([]byte("[ENTER YOUR NAME]: "))
	}

	// 3. Enforce connection limit
	if hub.ClientCount() >= maxClients {
		conn.Write([]byte("Server is full. Maximum 10 connections allowed. Try again later.\n"))
		return
	}

	// 4. Create client and start its dedicated writer goroutine
	client := server.NewClient(conn, name)
	go client.WritePump()

	// 5. Deliver full chat history BEFORE registering so history always arrives
	//    before any concurrent broadcasts from other clients.
	for _, msg := range history.GetAll() {
		client.Deliver(msg)
	}

	// 6. Register in the hub — from this point others can broadcast to this client
	hub.Register(client)

	// 7. Notify existing clients that this client joined
	joinMsg := fmt.Sprintf("%s has joined our chat...\n", name)
	hub.BroadcastExcept(joinMsg, client)
	hub.RepromptExcept(client) // re-show prompt to clients whose screen was interrupted

	// 8. Show the initial input prompt to the new client
	client.Deliver(client.Prompt())

	// 9. Main message loop
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // client disconnected or connection error
		}

		text := strings.TrimSpace(line)
		if text == "" {
			// Do not broadcast empty messages; just re-show the prompt
			client.Deliver(client.Prompt())
			continue
		}

		// Format: [YYYY-MM-DD HH:MM:SS][name]:message
		msg := fmt.Sprintf("[%s][%s]:%s\n", time.Now().Format("2006-01-02 15:04:05"), name, text)
		history.Add(msg)
		hub.BroadcastExcept(msg, client)
		hub.RepromptExcept(client) // re-show prompt to all other clients
		client.Deliver(client.Prompt())
	}

	// 10. Client disconnected — clean up
	hub.Unregister(client) // remove first so leave msg is not sent to this client

	leaveMsg := fmt.Sprintf("%s has left our chat...\n", name)
	history.Add(leaveMsg)
	hub.BroadcastAll(leaveMsg)
	hub.RepromptAll() // re-show prompt to remaining clients

	client.Close() // signal WritePump to exit
}

// GetWelcomeMessage reads the welcome ASCII art from disk.
// Falls back to a plain text message if the file cannot be read.
func GetWelcomeMessage() string {
	data, err := os.ReadFile("assets/welcome.txt")
	if err != nil {
		return "Welcome to TCP-Chat!\n[ENTER YOUR NAME]: \n"
	}
	return string(data)
}
