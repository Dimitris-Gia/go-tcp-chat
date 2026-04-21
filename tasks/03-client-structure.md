# Task 03: Client Data Structure and Management

## Status: ⏳ Not Started

## Objective
Create a data structure to represent connected clients and implement thread-safe management of the client list.

## Requirements
1. Define a Client struct with:
   - Connection (net.Conn)
   - Name (string)
   - Unique identifier
   
2. Create a ClientManager to handle:
   - Adding new clients
   - Removing clients
   - Getting list of all clients
   - Thread-safe operations using mutex or channels

3. Implement methods:
   - AddClient(client *Client)
   - RemoveClient(clientID)
   - GetAllClients() []*Client
   - GetClientCount() int

## TDD Steps
1. Write test for Client struct creation
2. Write test for adding client to manager
3. Write test for removing client from manager
4. Write test for concurrent access (goroutine safety)
5. Write test for max connection limit check

## Acceptance Criteria
- Client struct properly defined
- ClientManager handles concurrent access safely
- Tests pass for add/remove operations
- Can retrieve client count

## Files to Create/Modify
- `internal/client/client.go` (new)
- `internal/client/client_test.go` (new)
- `internal/client/manager.go` (new)
- `internal/client/manager_test.go` (new)

## Notes
- Use sync.Mutex for thread safety
- Consider using channels as alternative to mutex
- Client ID can be generated using time or counter
