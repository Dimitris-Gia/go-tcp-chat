# Task 09: Leave Notifications

## Status: ⏳ Not Started

## Objective
Notify remaining clients when a client leaves the chat.

## Requirements
1. Detect client disconnection (EOF, connection error)
2. Send leave notification to remaining clients
3. Format: `username has left our chat...`
4. Remove client from client manager
5. Add notification to message history
6. Clean up resources

## TDD Steps
1. Write test for leave notification format
2. Write test for detecting disconnection
3. Write test for removing client from manager
4. Write test for sending to remaining clients
5. Write test for resource cleanup

## Acceptance Criteria
- Disconnection detected properly (EOF, errors)
- Leave notification sent to remaining clients
- Client removed from manager
- Notification added to history
- Connection closed and resources cleaned up
- Tests verify all scenarios

## Files to Create/Modify
- `internal/notification/notification.go` (modify)
- `internal/notification/notification_test.go` (modify)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- Use defer for cleanup in connection handler
- Handle both graceful disconnect and connection errors
- Leave notification should not have timestamp prefix