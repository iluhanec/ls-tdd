<!-- 3ca7cdde-b294-45b9-81cf-0ef8ea90ec67 2d780c7f-0dbf-430b-b6b9-5e927e75e8a6 -->
# TDD Implementation Plan for Minimal ls Utility

## Overview

Implement a simplified version of the Linux `ls` command using Test-Driven Development. Focus on the 20% most popular functionality: basic listing, directory argument, long format (-l), and hidden files (-a).

## Most Popular Features (20% of ls functionality)

Based on common usage patterns, we'll implement:

1. **Basic listing** - List files in current directory (no flags)
2. **Directory argument** - List files in specified directory (`ls <path>`)
3. **Long format** - Detailed file information (`-l` flag)
4. **Hidden files** - Include dotfiles (`-a` flag)

## TDD Approach

Follow the Red-Green-Refactor cycle for each feature:

- **Red**: Write the simplest failing test
- **Green**: Write minimal code to make it pass
- **Refactor**: Improve design while keeping tests green
- **Reflect**: Review design and adjust plan if needed

## Implementation Steps

### Phase 1: Basic Infrastructure

**Step 1.1**: Simplest possible test

- Test: `ls` command exists and can be invoked
- Implementation: Create main function that prints empty output
- Reflection: Establish basic command structure

**Step 1.2**: Basic output test

- Test: `ls` outputs something (even if just a placeholder)
- Implementation: Print a single string
- Reflection: Establish output mechanism

### Phase 2: Basic Listing

**Step 2.1**: List current directory (simplest case)

- Test: `ls` lists at least one file from current directory
- Implementation: Read current directory, print first file name
- Reflection: Establish directory reading mechanism

**Step 2.2**: List all files in current directory

- Test: `ls` lists all non-hidden files in current directory
- Implementation: Read and print all non-hidden files
- Reflection: Consider file filtering logic

**Step 2.3**: Sort output

- Test: Files are listed in alphabetical order
- Implementation: Sort file list before printing
- Reflection: Establish sorting as default behavior

### Phase 3: Directory Argument

**Step 3.1**: Accept directory path argument

- Test: `ls <directory>` lists files in specified directory
- Implementation: Parse command-line argument, read that directory
- Reflection: Establish argument parsing structure

**Step 3.2**: Handle invalid directory

- Test: `ls <invalid-path>` returns appropriate error
- Implementation: Error handling for non-existent paths
- Reflection: Consider error message format

### Phase 4: Long Format (-l flag)

**Step 4.1**: Parse -l flag

- Test: `ls -l` triggers long format mode
- Implementation: Parse flag, set flag variable
- Reflection: Establish flag parsing mechanism

**Step 4.2**: Display file permissions

- Test: Long format shows file permissions (e.g., `-rw-r--r--`)
- Implementation: Get and format file mode
- Reflection: Consider permission formatting

**Step 4.3**: Display file size

- Test: Long format shows file size in bytes
- Implementation: Get file size, format as number
- Reflection: Consider size formatting

**Step 4.4**: Display modification time

- Test: Long format shows modification date/time
- Implementation: Get file modification time, format
- Reflection: Consider time format (similar to ls)

**Step 4.5**: Complete long format line

- Test: Long format shows all fields: permissions, links, owner, group, size, date, name
- Implementation: Combine all fields in proper format
- Reflection: Review complete output format

### Phase 5: Hidden Files (-a flag)

**Step 5.1**: Parse -a flag

- Test: `ls -a` includes hidden files in output
- Implementation: Parse -a flag, modify filtering logic
- Reflection: Consider flag combination (e.g., `ls -la`)

**Step 5.2**: Combine flags

- Test: `ls -la` shows hidden files in long format
- Implementation: Support multiple flags simultaneously
- Reflection: Review flag parsing design

## Design Considerations

### File Structure

- `main.go`: Entry point, command-line parsing, orchestration
- `ls.go`: Core ls logic (directory reading, file listing)
- `formatter.go`: Output formatting (basic and long format)
- `main_test.go`: Test file for main package

### Key Design Decisions to Make During Implementation

1. **Flag parsing**: Simple manual parsing vs. flag package (start simple, refactor if needed)
2. **File reading**: Use `os.ReadDir` or `filepath.Walk` (start with ReadDir)
3. **Output formatting**: String building vs. template (start with string building)
4. **Error handling**: Return errors vs. exit with code (start with return errors, handle in main)

### Testing Strategy

- Use Go's built-in `testing` package
- Create test directories and files as needed
- Test both success and error cases
- Use table-driven tests where appropriate

## Reflection Points

After each phase, consider:

- Is the current design extensible for new features?
- Are there any code smells or duplication?
- Can the code be simplified?
- Are edge cases handled?
- Should the design be refactored before adding next feature?

## Success Criteria

- All tests pass
- Code follows Go best practices
- Implementation handles common edge cases
- Design is clean and maintainable
- Documentation is clear