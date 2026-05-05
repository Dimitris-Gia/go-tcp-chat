# NetCat TCP Chat

Go application that recreates NetCat functionality in a Server-Client Architecture for group chat communication over TCP.

## Description

NetCat TCP Chat is a command-line group chat application that runs in server mode, listening for incoming TCP connections on a specified port. Multiple clients can connect simultaneously (up to 10) and communicate in real-time. The server broadcasts messages to all connected clients, maintains chat history, notifies users when others join or leave, and writes an audit log of all events.

## Authors

Dimitris Giannakakis, Panagiotis Ermidis, Eleutheria Manola

## Usage

### How to Run the Server

```bash
# Start server on default port (8989)
go run .
# Output: Listening on the port :8989

# Start server on custom port
go run . 2525
# Output: Listening on the port :2525

# Invalid usage (too many arguments)
go run . 2525 localhost
# Output: [USAGE]: ./TCPChat $port
```

### How to Connect as Client

```bash
nc localhost 8989
# or
nc localhost 2525
```

### Client Experience

1. Upon connection, client sees Linux penguin ASCII art
2. Client is prompted to enter their name: `[ENTER YOUR NAME]:`
3. After entering name, client receives full chat history
4. Client can send messages which are broadcast to all other clients
5. Client sees notifications when others join or leave
6. Messages are formatted: `[YYYY-MM-DD HH:MM:SS][username]:[message]`

## Features

- **TCP Server-Client Architecture**: 1-to-many relationship
- **Connection Management**: Maximum 10 concurrent connections
- **Name Registration**: Clients must provide non-empty names
- **Message Broadcasting**: Real-time message distribution to all clients
- **Chat History**: New clients receive all previous messages
- **Join/Leave Notifications**: Server informs clients of user activity
- **Timestamped Messages**: All messages include timestamp and sender name
- **Empty Message Filtering**: Empty messages are not broadcast
- **Graceful Disconnection**: Other clients remain connected when one leaves
- **Username Change**: Clients can rename themselves mid-session
- **Color Flags**: Messages can be colored using flag prefixes
- **Emote Flags**: Shorthand flags expand to ASCII emoticons
- **ANSI Line Clearing**: Incoming messages clear the prompt line to reduce visual interruption
- **Audit Logging**: All events written as JSON to `logs/audit.log`

## Commands

### Username Change

```
--UserNameChange: <new name>
```

Renames the client. All other clients are notified and the change is stored in history.

### Color Flags

Prefix your message with a color flag:

```
--red hello everyone
--green good morning
--yellow warning!
--blue just chilling
--magenta look at me
--cyan cool message
```

### Emote Flags

Send a standalone emote flag:

```
--shrug       →  \_(o_o)_/
--happy       →  ^_^
--sad         →  (｡•́︿•̀｡)
--wow         →  (⚆_⚆)
--heart       →  <3
--tableflip   →  (╯°□°）╯︵ ┻━┻
--unflip      →  ┬─┬ ノ( ゜-゜ノ)
--lenny       →  ( ͡° ͜ʖ ͡°)
--disapprove  →  ಠ_ಠ
--cry         →  T_T
--kiss        →  ( ˘ ³˘)♥
--weeping     →  (╥﹏╥)
--angry       →  ಠ益ಠ
--confused    →  (;一_一)
--party       →  └(＾ω＾)」
--sleepy      →  (-_-) zzz
```

## Implementation Details

### Algorithm

1. **Server Initialization**
   - Parse port from command-line arguments (default: 8989)
   - Validate port range (1-65535)
   - Initialize logger (`logs/audit.log`)
   - Start TCP listener on specified port
   - Accept incoming connections in loop

2. **Connection Handling**
   - Spawn goroutine for each client connection
   - Send welcome ASCII art from `assets/welcome.txt`
   - Prompt for and read client name (re-prompt on empty)
   - Enforce max 10 connection limit
   - Create `Client` with buffered send channel, start `WritePump` goroutine
   - Deliver chat history before registering in hub

3. **Message Flow**
   - Broadcast join notification to existing clients
   - Continuously read messages from client
   - Apply chat flags (color/emote) if present
   - Format messages with timestamp and sender name
   - Store in history, broadcast to all other clients
   - Filter out empty messages
   - Reprompt all other clients after each broadcast

4. **Disconnection Handling**
   - Detect client disconnection (EOF or error)
   - Unregister from hub
   - Broadcast leave notification to remaining clients
   - Store leave message in history
   - Close client (signals `WritePump` to exit)

5. **Concurrency Management**
   - Goroutines for concurrent client handling
   - `Hub` uses mutex to protect client map
   - `History` uses mutex to protect entries slice
   - `Client` uses `done` channel + `sync.Once` for safe close
   - `Client` uses mutex for name access only
   - `Logger` uses mutex for file writes

6. **Message Format**
   - Timestamp format: `[YYYY-MM-DD HH:MM:SS]`
   - Message format: `[timestamp][username]:[message]`
   - Join notification: `username has joined our chat...`
   - Leave notification: `username has left our chat...`
   - Rename notification: `oldname is now known as newname`

7. **Terminal UX**
   - Broadcasts prepend `\r\033[2K` to clear the recipient's current line
   - Prompt is appended immediately after the message in the same delivery
   - Reduces visual interruption when a message arrives mid-typing

### Welcome ASCII Art

```
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

## Audit Logging

All events are written as newline-delimited JSON to `logs/audit.log`:

```json
{"timestamp":"2024-01-01T12:00:00Z","level":"INFO","eventType":"ServerStarted","data":{"port":"8989"}}
{"timestamp":"2024-01-01T12:00:05Z","level":"INFO","eventType":"ClientJoined","data":{"clientName":"alice","ip":"127.0.0.1:54321","port":""}}
{"timestamp":"2024-01-01T12:00:10Z","level":"INFO","eventType":"MessageSent","data":{"content":"hello","sender":"alice"}}
{"timestamp":"2024-01-01T12:00:20Z","level":"INFO","eventType":"ClientDisconnected","data":{"clientName":"alice"}}
```

Event types: `ServerStarted`, `ServerStopped`, `ClientJoined`, `ClientDisconnected`, `MessageSent`, `NameChanged`, `Error`

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run tests with verbose output
go test -v ./...
```

### Manual Testing

1. Start the server in one terminal
2. Open multiple terminal windows (up to 10)
3. Connect using `nc localhost 8989` in each
4. Enter different names for each client
5. Send messages and verify broadcasting
6. Test color flags: `--red hello`
7. Test emote flags: `--shrug`
8. Test username change: `--UserNameChange: newname`
9. Disconnect clients and verify notifications
10. Test connection limit (11th connection should be rejected)

## Project Structure

```
net-cat/
├── main.go                          # Server entry point
├── main_test.go                     # Integration tests
├── assets/
│   └── welcome.txt                  # Linux penguin ASCII art
├── logs/
│   └── audit.log                    # Runtime audit log (JSON)
├── internal/
│   ├── connectionhandling/          # Client connection lifecycle
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── logging/                     # Audit logger
│   │   ├── events.go
│   │   ├── logger.go
│   │   └── logger_test.go
│   ├── parser/                      # CLI argument parsing
│   │   ├── parsing.go
│   │   └── parsing_test.go
│   └── server/                      # Core server types
│       ├── client.go
│       ├── client_test.go
│       ├── history.go
│       ├── history_test.go
│       ├── hub.go
│       └── hub_test.go
├── ai/
│   └── ai.txt                       # Conversation logs
├── tasks/                           # Implementation task tracking
├── agents.md
├── prd.md
└── README.md
```

## Allowed Packages

- `io` - Input/output operations
- `log` - Logging
- `os` - Operating system functionality
- `fmt` - Formatted I/O
- `net` - Network operations
- `sync` - Synchronization primitives
- `time` - Time operations
- `bufio` - Buffered I/O
- `errors` - Error handling
- `strings` - String manipulation
- `reflect` - Reflection
- `encoding/json` - JSON marshalling (for audit log)
- `path/filepath` - File path utilities (for audit log)

## Error Handling

- Invalid port number: display usage message
- Too many arguments: display usage message
- Connection errors: log and continue accepting new connections
- Client disconnection: clean up and notify others
- Maximum connections reached: reject new connection gracefully
- Invalid flag usage: return error message to sender only
- Logger failure: fatal on startup, logged on runtime errors
