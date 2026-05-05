package connectionhandling

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"net-cat/internal/logging"
	"net-cat/internal/server"
)

const maxClients = 10
const usernameChangePrefix = "--UserNameChange:"

var colorFlags = map[string]string{
	"--red":     "31",
	"--green":   "32",
	"--yellow":  "33",
	"--blue":    "34",
	"--magenta": "35",
	"--cyan":    "36",
}

var emoteFlags = map[string]string{
	"--shrug":      "\\_(o_o)_/",
	"--happy":      "^_^",
	"--sad":        "(｡•́︿•̀｡)",
	"--wow":        "(⚆_⚆)",
	"--heart":      "<3",
	"--tableflip":  "(╯°□°）╯︵ ┻━┻",
	"--unflip":     "┬─┬ ノ( ゜-゜ノ)",
	"--lenny":      "( ͡° ͜ʖ ͡°)",
	"--disapprove": "ಠ_ಠ",
	"--cry":        "T_T",
	"--kiss":       "( ˘ ³˘)♥",
	"--weeping":    "(╥﹏╥)",
	"--angry":      "ಠ益ಠ",
	"--confused":   "(;一_一)",
	"--party":      "└(＾ω＾)」",
	"--sleepy":     "(-_-) zzz",
}

// HandleConnection manages the full lifecycle of a single client connection.
// It must be run in its own goroutine.
func HandleConnection(conn net.Conn, hub *server.Hub, history *server.History, logger *logging.Logger) {
	defer conn.Close()

	// 1. Send welcome message (ASCII art + [ENTER YOUR NAME]: prompt)
	conn.Write([]byte(GetWelcomeMessage()))

	reader := bufio.NewReader(conn)

	// 2. Name registration loop — keep asking until a non-empty, unique name is given
	var name string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		name = strings.TrimSpace(line)
		if name == "" {
			conn.Write([]byte("[ENTER YOUR NAME]: "))
			continue
		}
		if hub.IsNameTaken(name) {
			conn.Write([]byte("Name already in use. [ENTER YOUR NAME]: "))
			continue
		}
		break
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

	// Log client join with IP and port information
	if logger != nil {
		clientAddr := conn.RemoteAddr().String()
		logger.LogEvent(logging.LevelInfo, logging.EventClientJoined, logging.ClientJoinedData(client.GetName(), clientAddr, ""))
	}

	// 7. Notify existing clients that this client joined
	joinMsg := fmt.Sprintf("%s has joined our chat...\n", client.GetName())
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

		if strings.HasPrefix(text, usernameChangePrefix) {
			newName := strings.TrimSpace(strings.TrimPrefix(text, usernameChangePrefix))
			if newName == "" {
				client.Deliver("\nInvalid name. Usage: --UserNameChange: <new name>\n")
				client.Deliver(client.Prompt())
				continue
			}
			if hub.IsNameTaken(newName) {
				client.Deliver("\nName already in use.\n")
				client.Deliver(client.Prompt())
				continue
			}

			oldName := client.GetName()
			client.SetName(newName)

			// Log name change
			if logger != nil {
				logger.LogEvent(logging.LevelInfo, logging.EventNameChanged, logging.NameChangedData(oldName, newName))
			}

			notice := fmt.Sprintf("%s is now known as %s\n", oldName, newName)
			history.Add(notice)
			hub.BroadcastExcept(notice, client)
			hub.RepromptExcept(client)
			client.Deliver("\n" + notice)
			client.Deliver(client.Prompt())
			continue
		}

		if transformed, err, handled := applyChatFlags(text); handled {
			if err != nil {
				client.Deliver("\n" + err.Error() + "\n")
				client.Deliver(client.Prompt())
				continue
			}
			text = transformed
		}

		// Format: [YYYY-MM-DD HH:MM:SS][name]:message
		msg := fmt.Sprintf("[%s][%s]:%s\n", time.Now().Format("2006-01-02 15:04:05"), client.GetName(), text)
		history.Add(msg)

		// Log message sent
		if logger != nil {
			logger.LogEvent(logging.LevelInfo, logging.EventMessageSent, logging.MessageSentData(client.GetName(), text))
		}

		hub.BroadcastExcept(msg, client)
		hub.RepromptExcept(client) // re-show prompt to all other clients
		client.Deliver(client.Prompt())
	}

	// 10. Client disconnected — clean up
	hub.Unregister(client) // remove first so leave msg is not sent to this client

	leaveMsg := fmt.Sprintf("%s has left our chat...\n", client.GetName())
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

func applyChatFlags(text string) (string, error, bool) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return "", nil, false
	}

	flag := strings.ToLower(parts[0])

	if ansiCode, ok := colorFlags[flag]; ok {
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid color usage. Example: --red hello"), true
		}
		content := strings.TrimSpace(strings.TrimPrefix(text, parts[0]))
		return fmt.Sprintf("\x1b[%sm%s\x1b[0m", ansiCode, content), nil, true
	}

	if emote, ok := emoteFlags[flag]; ok {
		return emote, nil, true
	}

	return "", nil, false
}
