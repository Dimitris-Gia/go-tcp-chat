# NetCat TCP Chat

Go application that recreates NetCat functionality in a Server-Client Architecture for group chat communication over TCP.

## Description

NetCat TCP Chat is a command-line group chat application that runs in server mode, listening for incoming TCP connections on a specified port. Multiple clients can connect simultaneously (up to 10) and communicate in real-time. The server broadcasts messages to all connected clients, maintains chat history, and notifies users when others join or leave.

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
# Connect to server using netcat
nc localhost 8989

# Or with custom port
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

## Implementation Details

### Algorithm

1. **Server Initialization**
   - Parse port from command-line arguments (default: 8989)
   - Validate port range (1-65535)
   - Start TCP listener on specified port
   - Accept incoming connections in loop

2. **Connection Handling**
   - Spawn goroutine for each client connection
   - Send welcome ASCII art (Linux penguin)
   - Prompt for and read client name
   - Validate non-empty name
   - Add client to active connections list

3. **Message Flow**
   - Send chat history to newly connected client
   - Broadcast join notification to existing clients
   - Continuously read messages from client
   - Format messages with timestamp and sender name
   - Broadcast to all other connected clients
   - Filter out empty messages

4. **Disconnection Handling**
   - Detect client disconnection (EOF or error)
   - Remove client from active connections
   - Broadcast leave notification to remaining clients
   - Close connection and clean up resources

5. **Concurrency Management**
   - Use goroutines for concurrent client handling
   - Use mutexes or channels to protect shared state
   - Shared state includes: client list, message history
   - Thread-safe broadcasting mechanism

6. **Message Format**
   - Timestamp format: `[YYYY-MM-DD HH:MM:SS]`
   - Message format: `[timestamp][username]:[message]`
   - Join notification: `username has joined our chat...`
   - Leave notification: `username has left our chat...`

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
6. Disconnect clients and verify notifications
7. Test connection limit (11th connection should be rejected)

## Project Structure

```
net-cat/
├── main.go                          # Server entry point
├── internal/
│   ├── connectionhandling/          # Client connection handling
│   │   └── handler.go
│   └── parser/                      # Command-line argument parsing
│       └── parsing.go
├── ai/                              # Conversation logs
│   └── ai.txt
├── tasks/                           # Implementation task tracking
├── prd.md                           # Product requirements
├── agents.md                        # Agent behavior guidelines
└── README.md                        # This file
```

## Allowed Packages

- `io` - Input/output operations
- `log` - Logging
- `os` - Operating system functionality
- `fmt` - Formatted I/O
- `net` - Network operations
- `sync` - Synchronization primitives (mutexes)
- `time` - Time operations
- `bufio` - Buffered I/O
- `errors` - Error handling
- `strings` - String manipulation
- `reflect` - Reflection

## Error Handling

- Invalid port number: Display usage message
- Too many arguments: Display usage message
- Connection errors: Log and continue accepting new connections
- Client disconnection: Clean up and notify others
- Maximum connections reached: Reject new connection gracefully

## Development Guidelines

- Follow TDD methodology (Red-Green-Refactor)
- Write tests before implementation
- Use only standard Go library packages
- Maintain conversation log in `ai/ai.txt`
- Return "ERROR" string for file format issues
- Keep code modular and well-organized
