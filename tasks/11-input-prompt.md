# Task 11: Input Prompt Handling

## Status: ⏳ Not Started

## Objective
Display input prompt for clients and handle message input properly.

## Requirements
1. Display prompt: `[YYYY-MM-DD HH:MM:SS][username]:`
2. Update prompt timestamp when other clients send messages
3. Maintain proper terminal display
4. Handle cursor positioning
5. Allow client to type while receiving messages

## TDD Steps
1. Write test for prompt format with timestamp
2. Write test for prompt update on incoming messages
3. Write test for concurrent input/output handling
4. Write test for proper message display
5. Write test for terminal formatting

## Acceptance Criteria
- Prompt shows current timestamp and username
- Prompt updates when messages arrive from others
- Client can type while receiving messages
- Messages display properly without interfering with input
- Terminal formatting maintained
- Tests verify prompt behavior

## Files to Create/Modify
- `internal/prompt/prompt.go` (new)
- `internal/prompt/prompt_test.go` (new)
- `internal/connectionhandling/handler.go` (modify)

## Notes
- This is complex terminal handling
- May need to use channels for coordinating input/output
- Consider simpler approach: just show prompt before each input
- Focus on core functionality first