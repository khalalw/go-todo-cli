package commands

import (
	"bytes"
	"go-todo-cli/internal/todo"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Set a test file name
	FileToWrite = "test_todos.json"

	// Run the tests
	code := m.Run()

	// Clean up the test file
	os.Remove(FileToWrite)

	// Exit with the test result code
	os.Exit(code)
}

func TestAddCommand(t *testing.T) {
	todos := &todo.Todos{}
	dueDate := time.Now().AddDate(0, 0, 1) // Tomorrow
	tags := []string{"work", "urgent"}
	AddCommand([]string{"Test task"}, &dueDate, todo.High, todos, tags)
	if len(*todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(*todos))
	}
	if (*todos)[0].Task != "Test task" {
		t.Errorf("Expected task 'Test task', got '%s'", (*todos)[0].Task)
	}
	if (*todos)[0].DueDate == nil || !(*todos)[0].DueDate.Equal(dueDate) {
		t.Errorf("Expected due date '%s', got '%v'", dueDate.Format("2006-01-02"), (*todos)[0].DueDate)
	}
	if (*todos)[0].Priority != todo.High {
		t.Errorf("Expected priority 'high', got '%s'", (*todos)[0].Priority)
	}
	if len((*todos)[0].Tags) != 2 || (*todos)[0].Tags[0] != "work" || (*todos)[0].Tags[1] != "urgent" {
		t.Errorf("Expected tags [work urgent], got %v", (*todos)[0].Tags)
	}
}

func TestCompleteCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task", Completed: false}}
	CompleteCommand([]string{"1"}, todos)
	if !(*todos)[0].Completed {
		t.Error("CompleteCommand failed to mark task as complete")
	}

	// Test invalid task number
	CompleteCommand([]string{"2"}, todos)
	if len(*todos) != 1 {
		t.Error("CompleteCommand should not add new tasks")
	}
}

func TestDeleteCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task"}}
	DeleteCommand([]string{"1"}, todos)
	if len(*todos) != 0 {
		t.Errorf("Expected 0 todos after deletion, got %d", len(*todos))
	}

	// Test invalid task number
	todos = &todo.Todos{{Task: "Test task"}}
	DeleteCommand([]string{"2"}, todos)
	if len(*todos) != 1 {
		t.Error("DeleteCommand should not remove tasks for invalid numbers")
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

func TestEditCommand(t *testing.T) {
	// Setup
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	tests := []struct {
		name             string
		input            string
		expectedTask     string
		expectedDueDate  string
		expectedPriority todo.Priority
		expectedTags     []string
	}{
		{
			name:             "Edit all fields",
			input:            "New task description\n2023-07-01\nhigh\nwork,urgent\n",
			expectedTask:     "New task description",
			expectedDueDate:  "2023-07-01",
			expectedPriority: todo.High,
			expectedTags:     []string{"work", "urgent"},
		},
		{
			name:             "Keep original values",
			input:            "\n\n\n\n",
			expectedTask:     "Original task",
			expectedDueDate:  "",
			expectedPriority: todo.Low,
			expectedTags:     []string{},
		},
		{
			name:             "Edit task only",
			input:            "Updated task\n\n\n\n",
			expectedTask:     "Updated task",
			expectedDueDate:  "",
			expectedPriority: todo.Low,
			expectedTags:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new todos slice for each test case
			todos := &todo.Todos{
				{Task: "Original task", DueDate: nil, Priority: todo.Low, Tags: []string{}},
			}

			// Create pipes for input and output
			inR, inW, _ := os.Pipe()
			outR, outW, _ := os.Pipe()
			os.Stdin = inR
			os.Stdout = outW

			// Write the test input to the pipe
			go func() {
				inW.Write([]byte(tt.input))
				inW.Close()
			}()

			// Run the command
			EditCommand(1, todos)

			// Close the write end of the output pipe
			outW.Close()

			// Read the output
			var output bytes.Buffer
			io.Copy(&output, outR)

			// Check the result
			if (*todos)[0].Task != tt.expectedTask {
				t.Errorf("Expected task '%s', got '%s'", tt.expectedTask, (*todos)[0].Task)
			}

			if tt.expectedDueDate != "" {
				expectedDate, _ := time.Parse("2006-01-02", tt.expectedDueDate)
				if (*todos)[0].DueDate == nil || !(*todos)[0].DueDate.Equal(expectedDate) {
					t.Errorf("Expected due date '%s', got '%v'", tt.expectedDueDate, (*todos)[0].DueDate)
				}
			} else if (*todos)[0].DueDate != nil {
				t.Errorf("Expected no due date, got '%v'", (*todos)[0].DueDate)
			}

			if (*todos)[0].Priority != tt.expectedPriority {
				t.Errorf("Expected priority '%v', got '%v'", tt.expectedPriority, (*todos)[0].Priority)
			}

			if !stringSlicesEqual((*todos)[0].Tags, tt.expectedTags) {
				t.Errorf("Expected tags %v, got %v", tt.expectedTags, (*todos)[0].Tags)
			}

			// Check the output if needed
			if !strings.Contains(output.String(), "Task updated successfully") {
				t.Errorf("Expected output to contain 'Task updated successfully', got '%s'", output.String())
			}
		})
	}
}

func TestAddTagCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task", Tags: []string{"existing"}}}

	// Test adding a new tag
	AddTagCommand([]string{"1", "newtag"}, todos)
	if len((*todos)[0].Tags) != 2 || (*todos)[0].Tags[1] != "newtag" {
		t.Errorf("Expected 2 tags with 'newtag' added, got %v", (*todos)[0].Tags)
	}

	// Test adding an existing tag
	AddTagCommand([]string{"1", "existing"}, todos)
	if len((*todos)[0].Tags) != 2 {
		t.Errorf("Expected no change when adding existing tag, got %v", (*todos)[0].Tags)
	}

	// Test invalid task number
	AddTagCommand([]string{"999", "tag"}, todos)
	if len((*todos)[0].Tags) != 2 {
		t.Errorf("Expected no change for invalid task number, got %v", (*todos)[0].Tags)
	}
}

func TestRemoveTagCommand(t *testing.T) {
	todos := &todo.Todos{{Task: "Test task", Tags: []string{"tag1", "tag2"}}}

	// Test removing an existing tag
	RemoveTagCommand([]string{"1", "tag1"}, todos)
	if len((*todos)[0].Tags) != 1 || (*todos)[0].Tags[0] != "tag2" {
		t.Errorf("Expected 1 tag 'tag2' remaining, got %v", (*todos)[0].Tags)
	}

	// Test removing a non-existent tag
	RemoveTagCommand([]string{"1", "nonexistent"}, todos)
	if len((*todos)[0].Tags) != 1 {
		t.Errorf("Expected no change when removing non-existent tag, got %v", (*todos)[0].Tags)
	}

	// Test invalid task number
	RemoveTagCommand([]string{"999", "tag2"}, todos)
	if len((*todos)[0].Tags) != 1 {
		t.Errorf("Expected no change for invalid task number, got %v", (*todos)[0].Tags)
	}
}

func TestFilterByTagCommand(t *testing.T) {
	todos := &todo.Todos{
		{Task: "Task 1", Tags: []string{"work", "urgent"}},
		{Task: "Task 2", Tags: []string{"personal"}},
		{Task: "Task 3", Tags: []string{"work"}},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	FilterByTagCommand([]string{"work"}, todos)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "Task 1") || !strings.Contains(output, "Task 3") || strings.Contains(output, "Task 2") {
		t.Errorf("FilterByTagCommand failed to filter tasks correctly")
	}

	// Test non-existent tag
	r, w, _ = os.Pipe()
	os.Stdout = w

	FilterByTagCommand([]string{"nonexistent"}, todos)

	w.Close()
	os.Stdout = old

	buf.Reset()
	io.Copy(&buf, r)

	output = buf.String()
	if !strings.Contains(output, "No tasks found with tag 'nonexistent'") {
		t.Errorf("FilterByTagCommand failed to handle non-existent tag correctly")
	}
}

func TestSearchCommand(t *testing.T) {
	todoList := &todo.Todos{}
	now := time.Now()
	todoList.Add("Buy groceries", &now, todo.Medium, []string{"shopping"})
	todoList.Add("Finish project", nil, todo.High, []string{"work"})
	todoList.Add("Call mom", nil, todo.Low, []string{"personal"})

	testCases := []struct {
		args     []string
		expected string
	}{
		{[]string{"groceries"}, "Buy groceries"},
		{[]string{"work"}, "Finish project"},
		{[]string{"personal"}, "Call mom"},
		{[]string{"nonexistent"}, "No tasks found matching 'nonexistent'"},
	}

	for _, tc := range testCases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			SearchCommand(tc.args, todoList)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if !strings.Contains(output, tc.expected) {
				t.Errorf("Expected output to contain '%s', but got:\n%s", tc.expected, output)
			}
		})
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
