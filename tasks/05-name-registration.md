# Task 05: Name Registration

## Status: ⏳ Not Started

## Objective
Prompt client for name and validate non-empty input.

## Requirements
1. Send prompt: `[ENTER YOUR NAME]:`
2. Read name input from client
3. Validate name is non-empty (trim whitespace)
4. Re-prompt if name is empty
5. Store name in Client struct

## TDD Steps
1. Write test for prompt message format
2. Write test for reading name from connection
3. Write test for empty name validation
4. Write test for whitespace-only name rejection
5. Write test for successful name registration

## Acceptance Criteria
- Prompt is displayed correctly
- Empty names are rejected
- Whitespace-only names are rejected
- Valid names are accepted and stored
- Client can retry if name is invalid
- Tests cover all validation cases

## Files to Create/Modify
- `internal/message/prompt.go` (new)
- `internal/message/prompt_test.go` (new)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- Use bufio.Reader to read input line by line
- Use strings.TrimSpace to clean input
- Consider timeout for name input (optional)
- Handle connection errors during read
