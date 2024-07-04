package commands

import (
	"bytes"
	"go-todo-cli/internal/todo"
	"io"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Check if todos.json exists and rename it if it does
	originalExists := false
	if _, err := os.Stat(FileToWrite); err == nil {
		originalExists = true
		os.Rename(FileToWrite, FileToWrite+".backup")
	}

	// Run the tests
	code := m.Run()

	// Clean up the test file
	os.Remove(FileToWrite)

	// Restore the original file if it existed
	if originalExists {
		os.Rename(FileToWrite+".backup", FileToWrite)
	}

	// Exit with the test result code
	os.Exit(code)
}

func TestAddCommand(t *testing.T) {
	todos := &todo.Todos{}
	dueDate := time.Now().AddDate(0, 0, 1) // Tomorrow
	AddCommand([]string{"Test task"}, &dueDate, todos)
	if len(*todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(*todos))
	}
	if (*todos)[0].Task != "Test task" {
		t.Errorf("Expected task 'Test task', got '%s'", (*todos)[0].Task)
	}
	if (*todos)[0].DueDate == nil || !(*todos)[0].DueDate.Equal(dueDate) {
		t.Errorf("Expected due date '%s', got '%v'", dueDate.Format("2006-01-02"), (*todos)[0].DueDate)
	}
}

func TestCompleteCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task", Completed: false}}
	CompleteCommand([]string{"1"}, todos)
	if !(*todos)[0].Completed {
		t.Error("CompleteCommand failed to mark task as complete")
	}
}

func TestDeleteCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task"}}
	DeleteCommand([]string{"1"}, todos)
	if len(*todos) != 0 {
		t.Errorf("Expected 0 todos after deletion, got %d", len(*todos))
	}
}

func TestListCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task"}}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ListCommand(todos)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if !bytes.Contains(buf.Bytes(), []byte("Test task")) {
		t.Error("ListCommand failed to list tasks")
	}
}

func TestClearTasksCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task"}}
	ClearTasksCommand(todos)
	if len(*todos) != 0 {
		t.Errorf("Expected 0 todos after clearing, got %d", len(*todos))
	}
}

func TestParseIndex(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"1", 0},
		{"5", 4},
		{"0", -1},
		{"-1", -1},
		{"abc", -1},
	}

	for _, tc := range testCases {
		result := parseIndex(tc.input)
		if result != tc.expected {
			t.Errorf("parseIndex(%s): expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}

func TestSaveTodoList(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task"}}
	saveTodoList(todos)

	// Check if file was created
	if _, err := os.Stat(FileToWrite); os.IsNotExist(err) {
		t.Error("saveTodoList did not create a file")
	}

	// Verify file contents
	data, err := os.ReadFile(FileToWrite)
	if err != nil {
		t.Errorf("Error reading saved file: %v", err)
	}
	if !bytes.Contains(data, []byte("Test task")) {
		t.Error("Saved file does not contain expected content")
	}
}
