# Task 06: Message History Storage

## Status: ⏳ Not Started

## Objective
Store all chat messages in history and send history to newly connected clients.

## Requirements
1. Create message history storage (slice of strings)
2. Thread-safe append to history
3. Thread-safe retrieval of history
4. Send full history to new client after name registration
5. Format: `[YYYY-MM-DD HH:MM:SS][username]:[message]`

## TDD Steps
1. Write test for adding message to history
2. Write test for retrieving all messages
3. Write test for concurrent access to history
4. Write test for sending history to connection
5. Write test for timestamp format

## Acceptance Criteria
- Messages stored with correct timestamp format
- History accessible in thread-safe manner
- New clients receive all previous messages
- Empty history handled correctly
- Tests verify concurrent access safety

## Files to Create/Modify
- `internal/history/history.go` (new)
- `internal/history/history_test.go` (new)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- Use sync.Mutex for thread safety
- Timestamp format: time.Format("2006-01-02 15:04:05")
- Consider memory limits for very long histories (optional)
