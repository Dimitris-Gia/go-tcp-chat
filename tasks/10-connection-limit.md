# Task 10: Connection Limit (Max 10)

## Status: ⏳ Not Started

## Objective
Enforce maximum of 10 concurrent client connections.

## Requirements
1. Check connection count before accepting new client
2. Reject connection if limit reached (10 clients)
3. Send rejection message to client before closing
4. Log rejection attempt
5. Continue accepting other connections

## TDD Steps
1. Write test for connection count tracking
2. Write test for accepting connection under limit
3. Write test for rejecting connection at limit
4. Write test for rejection message
5. Write test for continuing to accept after rejection

## Acceptance Criteria
- Maximum 10 clients can connect simultaneously
- 11th connection attempt is rejected gracefully
- Rejection message sent to client
- Server continues accepting new connections after rejection
- Connection count accurate after clients leave
- Tests verify limit enforcement

## Files to Create/Modify
- `internal/client/manager.go` (modify)
- `internal/client/manager_test.go` (modify)
- `main.go` (modify)

## Notes
- Check limit before calling HandleConnection
- Send polite rejection message before closing
- Consider "Server full, please try again later"
- Ensure count decreases when clients disconnect