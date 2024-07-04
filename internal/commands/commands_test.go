package commands

import (
	"bytes"
	"go-todo-cli/pkg/todo"
	"os"
	"strings"
	"testing"
)

// Utility function to capture the output of a function call
func captureOutput(f func()) string {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	return buf.String()
}

func TestAddCommand(t *testing.T) {
	todoList := &todo.Todos{}
	addCommand([]string{"Test task"}, todoList)

	if len(*todoList) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(*todoList))
	}

	if (*todoList)[0].Task != "Test task" {
		t.Errorf("Expected 'Test task', got %s", (*todoList)[0].Task)
	}
}

func TestCompleteCommand(t *testing.T) {
	todoList := &todo.Todos{
		{Task: "Task 1", Completed: false},
	}

	completeCommand([]string{"1"}, todoList)

	if !(*todoList)[0].Completed {
		t.Errorf("Expected task to be marked as complete")
	}
}

func TestDeleteCommand(t *testing.T) {
	todoList := &todo.Todos{
		{Task: "Task 1", Completed: false},
	}

	deleteCommand([]string{"1"}, todoList)

	if len(*todoList) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(*todoList))
	}
}

func TestClearTasksCommand(t *testing.T) {
	todoList := &todo.Todos{
		{Task: "Task 1", Completed: false},
		{Task: "Task 2", Completed: false},
	}

	clearTasksCommand(todoList)

	if len(*todoList) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(*todoList))
	}
}

func TestListCommand(t *testing.T) {
	todoList := &todo.Todos{
		{Task: "Task A", Completed: true},
		{Task: "Task B", Completed: false},
	}

	output := captureOutput(func() {
		listCommand(todoList)
	})

	expectedOutput := "1. [x] Task A\n2. [ ] Task B\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestExitCommand(t *testing.T) {
	todoList := &todo.Todos{
		{Task: "Task 1", Completed: false},
	}
	defer func() {
		if r := recover(); r != nil {
			if r != "os.Exit called" {
				panic(r)
			}
		}
	}()
	// Mock os.Exit to prevent it from terminating the test
	exitFunc = func(code int) {
		panic("os.Exit called")
	}
	defer func() { exitFunc = os.Exit }() // Restore original function

	output := captureOutput(func() {
		exitCommand(todoList)
	})

	expectedOutput := "Exiting the TODO CLI.\n"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain:\n%s\nGot:\n%s", expectedOutput, output)
	}
}

func TestParseIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 0},
		{"2", 1},
		{"a", -1},
		{"-1", -1},
		{" 2 ", 1},
		{"", -1},
	}

	for _, test := range tests {
		result := parseIndex(test.input)
		if result != test.expected {
			t.Errorf("parseIndex(%q) = %d; expected %d", test.input, result, test.expected)
		}
	}
}

func TestPrintHelp(t *testing.T) {
	expectedOutput := `Commands:
  add <task>           - Add a new task
  complete <task_number> - Mark a task as complete
  delete <task_number> - Delete a task
  list                 - List all tasks
  clear-tasks          - Clear all tasks
  help                 - Show this help message
  exit                 - Exit the TODO CLI
`
	output := captureOutput(func() {
		printHelp()
	})

	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
