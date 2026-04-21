# Task 04: Welcome Message and ASCII Art

## Status: ⏳ Not Started

## Objective
Display Linux penguin ASCII art welcome message when client connects.

## Requirements
1. Create welcome message with ASCII art:
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
```

2. Send welcome message to client upon connection
3. Handle write errors gracefully

## TDD Steps
1. Write test to verify welcome message content
2. Write test to verify message is sent to connection
3. Implement GetWelcomeMessage() function
4. Implement SendWelcome(conn net.Conn) function

## Acceptance Criteria
- Welcome message contains correct ASCII art
- Message is sent immediately upon connection
- Function handles connection write errors
- Tests verify message content and delivery

## Files to Create/Modify
- `internal/message/welcome.go` (new)
- `internal/message/welcome_test.go` (new)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- Store ASCII art as constant string
- Use fmt.Fprint or io.WriteString to send to connection
- Consider creating a message package for all message-related functions
