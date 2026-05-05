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
- Color and emote flags for message formatting
- Username change command
- Audit logging to file

### Constraints
- Language: Go
- Dependencies: standard library only (io, log, os, fmt, net, sync, time, bufio, errors, strings, reflect, encoding/json, path/filepath)
- Development style: TDD (Red-Green-Refactor)
- Conversation log must be maintained in `ai/ai.txt`
- Maximum 10 concurrent connections

### Completed Implementation ✅

1. **Port Parsing and Validation**
   - `GetPortNumber` validates port argument
   - Default port: 8989
   - Valid port range: 1-65535
   - Usage message: `[USAGE]: ./TCPChat $port`

2. **TCP Server**
   - Listens on specified port
   - Accepts multiple concurrent connections
   - Spawns goroutine per connection
   - Logs server start event

3. **Client Structure**
   - `Client` struct with buffered send channel (256 messages)
   - `WritePump` goroutine per client for non-blocking writes
   - Thread-safe `Deliver` using `done` channel + `sync.Once` (no mutex on send)
   - Thread-safe `GetName`/`SetName` via mutex
   - Idempotent `Close` via `sync.Once`
   - Formatted prompt: `[YYYY-MM-DD HH:MM:SS][name]:`

4. **Hub (Connection Manager)**
   - Thread-safe client registry using mutex
   - `Register`/`Unregister`/`ClientCount`
   - `BroadcastAll` — delivers to all clients with ANSI line clear + prompt
   - `BroadcastExcept` — delivers to all except sender with ANSI line clear + prompt
   - `RepromptAll` / `RepromptExcept` — redraws prompt with ANSI line clear

5. **Message History**
   - Thread-safe append-only slice
   - `Add` / `GetAll` (returns a copy)
   - Delivered to new clients before registration

6. **Connection Handling**
   - Welcome ASCII art from `assets/welcome.txt` (fallback to plain text)
   - Name registration loop — re-prompts on empty name
   - Connection limit enforced (max 10) after name registration
   - Chat history delivered before hub registration
   - Join notification broadcast to existing clients
   - Main message loop: format, store, broadcast
   - Empty messages filtered (not broadcast)
   - Leave notification on disconnect

7. **Username Change Command**
   - Command: `--UserNameChange: <new name>`
   - Validates non-empty new name
   - Notifies all clients and updates history
   - Updates sender's prompt immediately

8. **Chat Flags**
   - Color flags: `--red`, `--green`, `--yellow`, `--blue`, `--magenta`, `--cyan`
   - Emote flags: `--shrug`, `--happy`, `--sad`, `--wow`, `--heart`, `--tableflip`, `--unflip`, `--lenny`, `--disapprove`, `--cry`, `--kiss`, `--weeping`, `--angry`, `--confused`, `--party`, `--sleepy`
   - Invalid flag usage returns an error message to the sender

9. **Audit Logging**
   - JSON-formatted log entries written to `logs/audit.log`
   - Thread-safe via mutex
   - Events: `ServerStarted`, `ServerStopped`, `ClientJoined`, `ClientDisconnected`, `MessageSent`, `NameChanged`, `Error`
   - Each entry includes: timestamp (RFC3339), level, eventType, data map

10. **Terminal UX**
    - ANSI escape codes (`\r\033[2K`) clear the current line before delivering messages/prompts to other clients
    - Prevents incoming messages from visually corrupting the prompt line

### Functional Requirements

1. **TCP Server** ✅
   - ✅ Listen on configurable port (default 8989)
   - ✅ Accept multiple concurrent connections
   - ✅ Limit connections to maximum 10
   - ✅ Handle connection errors gracefully

2. **Client Connection Flow** ✅
   - ✅ Display Linux penguin ASCII art on connect
   - ✅ Prompt for client name: `[ENTER YOUR NAME]:`
   - ✅ Validate non-empty name
   - ✅ Send chat history to new client
   - ✅ Broadcast join notification to existing clients

3. **Message Handling** ✅
   - ✅ Read messages from clients continuously
   - ✅ Format: `[YYYY-MM-DD HH:MM:SS][client.name]:[client.message]`
   - ✅ Broadcast messages to all other connected clients
   - ✅ Do not broadcast empty messages
   - ✅ Store message history for new clients

4. **Client Disconnection** ✅
   - ✅ Detect client disconnection
   - ✅ Broadcast leave notification to remaining clients
   - ✅ Clean up client resources
   - ✅ Other clients remain connected

5. **Concurrency Management** ✅
   - ✅ Goroutines for concurrent client handling
   - ✅ Channels + mutexes for thread-safe operations
   - ✅ Protected shared state (client list, message history, logger)

6. **Input Prompt** ✅
   - ✅ Display prompt: `[YYYY-MM-DD HH:MM:SS][client.name]:`
   - ✅ Reprompt other clients after each broadcast
   - ✅ ANSI line clearing to reduce visual interruption

### Quality Requirements
- Tests cover happy-path and failure-path behavior
- Code is modular (separate packages: connectionhandling, server, parser, logging)
- Only standard Go library packages used
- Proper concurrent programming practices
- Proper error handling on both server and client side

### Success Criteria ✅
- ✅ Server starts and listens on specified port
- ✅ Multiple clients can connect simultaneously (up to 10)
- ✅ Clients see welcome ASCII art and name prompt
- ✅ New clients receive full chat history
- ✅ Messages are broadcast to all clients with proper format
- ✅ Join/leave notifications work correctly
- ✅ Empty messages are not broadcast
- ✅ Server handles client disconnections gracefully
- ✅ Connection limit enforced (max 10)
- ✅ All tests pass
- ✅ Audit log written to `logs/audit.log`

### Legend
- ✅ Completed
- 🚧 In Progress
- ⏳ Not Started
