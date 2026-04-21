# Task 07: Message Broadcasting

## Status: ⏳ Not Started

## Objective
Broadcast messages from one client to all other connected clients.

## Requirements
1. Read messages continuously from client
2. Format message with timestamp and username
3. Broadcast to all OTHER clients (not sender)
4. Filter out empty messages
5. Add message to history
6. Handle client disconnection during read

## TDD Steps
1. Write test for message formatting
2. Write test for broadcasting to multiple clients
3. Write test for filtering empty messages
4. Write test for not sending to sender
5. Write test for handling disconnected clients during broadcast

## Acceptance Criteria
- Messages formatted correctly: `[YYYY-MM-DD HH:MM:SS][username]:[message]`
- Empty messages not broadcast
- All clients except sender receive message
- Message added to history
- Disconnected clients handled gracefully
- Tests verify all scenarios

## Files to Create/Modify
- `internal/broadcast/broadcast.go` (new)
- `internal/broadcast/broadcast_test.go` (new)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- Use bufio.Scanner or bufio.Reader for reading lines
- Trim whitespace before checking if empty
- Handle write errors to individual clients
- Consider using channels for message queue (optional)