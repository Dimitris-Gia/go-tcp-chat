# Product Requirements Document (PRD)
## NetCat TCP Chat Application

### Overview
NetCat TCP Chat is a Go application that recreates the NetCat functionality in a Server-Client Architecture. It runs in server mode on a specified port, listening for incoming TCP connections, and allows multiple clients to connect and participate in a group chat.

### Scope
- TCP server listening on configurable port (default: 8989)
- Support multiple concurrent client connections (max 10)
- Client name registration requirement
- Real-time message broadcasting to all connected clients
- Message history for new clients
- Join/leave notifications
- Timestamped messages with sender identification

### Constraints
- Language: Go
- Dependencies: standard library only (io, log, os, fmt, net, sync, time, bufio, errors, strings, reflect)
- Development style: TDD (Red-Green-Refactor)
- Conversation log must be maintained in `ai/ai.txt`
- Maximum 10 concurrent connections

### Current Baseline

#### Completed Implementation ✅
1. **Port Parsing and Validation**
   - `GetPortNumber` function validates port argument
   - Default port: 8989
   - Valid port range: 1-65535
   - Usage message: `[USAGE]: ./TCPChat $port`

2. **Basic TCP Server**
   - Server listens on specified port
   - Accepts incoming connections
   - Spawns goroutine per connection
   - Basic connection handler structure

#### In Progress 🚧
3. **Connection Handling**
   - HandleConnection function skeleton exists
   - Needs implementation for welcome message, name registration, message handling

### Functional Requirements

1. **TCP Server** 🚧
   - ✅ Listen on configurable port (default 8989)
   - ✅ Accept multiple concurrent connections
   - ⏳ Limit connections to maximum 10
   - ⏳ Handle connection errors gracefully

2. **Client Connection Flow** ⏳
   - ⏳ Display Linux penguin ASCII art on connect
   - ⏳ Prompt for client name: `[ENTER YOUR NAME]:`
   - ⏳ Validate non-empty name
   - ⏳ Send chat history to new client
   - ⏳ Broadcast join notification to existing clients

3. **Message Handling** ⏳
   - ⏳ Read messages from clients continuously
   - ⏳ Format: `[YYYY-MM-DD HH:MM:SS][client.name]:[client.message]`
   - ⏳ Broadcast messages to all other connected clients
   - ⏳ Do not broadcast empty messages
   - ⏳ Store message history for new clients

4. **Client Disconnection** ⏳
   - ⏳ Detect client disconnection
   - ⏳ Broadcast leave notification to remaining clients
   - ⏳ Clean up client resources
   - ⏳ Other clients remain connected

5. **Concurrency Management** ⏳
   - ⏳ Use goroutines for concurrent client handling
   - ⏳ Use channels or mutexes for thread-safe operations
   - ⏳ Protect shared state (client list, message history)

6. **Input Prompt** ⏳
   - ⏳ Display prompt: `[YYYY-MM-DD HH:MM:SS][client.name]:`
   - ⏳ Update prompt when other clients send messages
   - ⏳ Maintain proper terminal display

### Quality Requirements
- Tests should cover happy-path and failure-path behavior
- Code should remain modular (separate packages for connection handling, message broadcasting, etc.)
- Use only standard Go packages
- Follow good practices for concurrent programming
- Proper error handling on both server and client side
- Unit tests recommended for server and client logic

### Success Criteria
- ⏳ Server starts and listens on specified port
- ⏳ Multiple clients can connect simultaneously (up to 10)
- ⏳ Clients see welcome ASCII art and name prompt
- ⏳ New clients receive full chat history
- ⏳ Messages are broadcast to all clients with proper format
- ⏳ Join/leave notifications work correctly
- ⏳ Empty messages are not broadcast
- ⏳ Server handles client disconnections gracefully
- ⏳ Connection limit enforced (max 10)
- ⏳ All tests pass
- ⏳ Documentation reflects actual implementation

### Legend
- ✅ Completed
- 🚧 In Progress
- ⏳ Not Started
