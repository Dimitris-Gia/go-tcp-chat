# Project: NET-CAT (Go)

## Agent Behavior Guidelines

1. **No Code Suggestions Unless Requested:**
   - **Do not suggest code or provide direct code solutions** unless I explicitly ask for help with a specific part of the code.
   - **Provide guidance, explanations, and clarifications** on concepts, best practices, or general programming advice when asked, but do not write code unless instructed.

2. **Test-Driven Development (TDD) Methodology:**
   - Follow the **TDD process** for the project: I will write a failing test first, then write the code to pass it. 
   - Encourage a focus on **writing tests first** and ensure that any functionality or feature we discuss is accompanied by a test case.
   - Make sure that we follow the **Red-Green-Refactor** cycle in TDD:
     - **Red**: Write a test that fails (since the code doesn’t exist yet).
     - **Green**: Write the minimum code necessary to pass the test.
     - **Refactor**: Clean up the code while ensuring the test still passes.

3. **Conversation Logging:**
   - Keep track of all conversations in a file called **ai.txt**, inside a folder named **ai**.
   - The log should contain the date and the model of AI used in the conversation. It should be updated automatically with each new interaction.

4. **No External Libraries or Code Outside the Standard Go Library:**
   - The project should be completed using only the **standard Go packages**. There should be no reliance on external libraries.
   
5. **Error Handling:**
   - If any errors occur, or there are issues with file formats, return the string **ERROR**.
