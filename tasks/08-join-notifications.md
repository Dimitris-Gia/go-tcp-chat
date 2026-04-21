# Task 08: Join Notifications

## Status: ⏳ Not Started

## Objective
Notify existing clients when a new client joins the chat.

## Requirements
1. Send join notification to all existing clients
2. Format: `username has joined our chat...`
3. Send notification after name registration and history delivery
4. Do not send notification to the joining client
5. Add notification to message history

## TDD Steps
1. Write test for join notification format
2. Write test for sending to existing clients only
3. Write test for adding notification to history
4. Write test for timing (after name registration)
5. Write test for multiple clients joining

## Acceptance Criteria
- Notification format correct: `username has joined our chat...`
- Only existing clients receive notification
- Joining client does not receive their own join notification
- Notification added to history for future clients
- Tests verify correct behavior

## Files to Create/Modify
- `internal/notification/notification.go` (new)
- `internal/notification/notification_test.go` (new)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- Join notification should not have timestamp prefix
- Consider it as a system message vs user message
- Ensure notification is sent after client is added to manager