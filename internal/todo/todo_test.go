package todo

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	todos := &Todos{}
	dueDate := time.Now().AddDate(0, 0, 1) // Tomorrow
	tags := []string{"work", "urgent"}
	todos.Add("Test task", &dueDate, High, tags)
	if len(*todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(*todos))
	}
	if (*todos)[0].Task != "Test task" {
		t.Errorf("Expected 'Test task', got '%s'", (*todos)[0].Task)
	}
	if (*todos)[0].DueDate == nil || !(*todos)[0].DueDate.Equal(dueDate) {
		t.Errorf("Expected due date '%s', got '%v'", dueDate.Format("2006-01-02"), (*todos)[0].DueDate)
	}
	if (*todos)[0].Priority != High {
		t.Errorf("Expected priority 'high', got '%s'", (*todos)[0].Priority)
	}
	if len((*todos)[0].Tags) != 2 || (*todos)[0].Tags[0] != "work" || (*todos)[0].Tags[1] != "urgent" {
		t.Errorf("Expected tags [work urgent], got %v", (*todos)[0].Tags)
	}
}

func TestComplete(t *testing.T) {
	todos := &Todos{{Task: "Test task", Completed: false}}
	err := todos.Complete(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !(*todos)[0].Completed {
		t.Error("Expected task to be completed")
	}
}

func TestDelete(t *testing.T) {
	todos := &Todos{{Task: "Test task"}}
	err := todos.Delete(0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(*todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(*todos))
	}
}

func TestSaveAndLoad(t *testing.T) {
	todos := &Todos{{Task: "Test task"}}
	filename := "test_todos.json"

	err := todos.Save(filename)
	if err != nil {
		t.Errorf("Error saving todos: %v", err)
	}

	loadedTodos := &Todos{}
	err = loadedTodos.Load(filename)
	if err != nil {
		t.Errorf("Error loading todos: %v", err)
	}

	if len(*loadedTodos) != 1 || (*loadedTodos)[0].Task != "Test task" {
		t.Error("Loaded todos do not match saved todos")
	}

	os.Remove(filename) // Clean up
}

func TestParsePriority(t *testing.T) {
	testCases := []struct {
		input    string
		expected Priority
		hasError bool
	}{
		{"low", Low, false},
		{"medium", Medium, false},
		{"high", High, false},
		{"LOW", Low, false},
		{"MEDIUM", Medium, false},
		{"HIGH", High, false},
		{"invalid", Low, true},
	}

	for _, tc := range testCases {
		result, err := ParsePriority(tc.input)
		if tc.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", tc.input)
		}
		if !tc.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
		}
		if result != tc.expected {
			t.Errorf("For input '%s', expected %v, but got %v", tc.input, tc.expected, result)
		}
	}
}

func TestPrint(t *testing.T) {
	todos := &Todos{
		{Task: "Task 1", Completed: false, Priority: High, Tags: []string{"urgent", "work"}},
		{Task: "Task 2", Completed: true, Priority: Low, Tags: []string{"personal"}},
	}
	dueDate := time.Now().AddDate(0, 0, 1)
	todos.Add("Task 3", &dueDate, Medium, []string{"project"})

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print(todos)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	expectedStrings := []string{"Task 1", "Task 2", "Task 3", "High", "Low", "Medium", "Pending", "Done", dueDate.Format("2006-01-02"), "urgent, work", "personal", "project"}

	for _, s := range expectedStrings {
		if !strings.Contains(output, s) {
			t.Errorf("Expected output to contain '%s', but it didn't", s)
		}
	}
}

func TestTodosOutOfRange(t *testing.T) {
	todos := &Todos{{Task: "Test task"}}

	// Test Complete out of range
	err := todos.Complete(-1)
	if err == nil {
		t.Error("Expected error for out of range index in Complete, but got none")
	}
	err = todos.Complete(1)
	if err == nil {
		t.Error("Expected error for out of range index in Complete, but got none")
	}

	// Test Delete out of range
	err = todos.Delete(-1)
	if err == nil {
		t.Error("Expected error for out of range index in Delete, but got none")
	}
	err = todos.Delete(1)
	if err == nil {
		t.Error("Expected error for out of range index in Delete, but got none")
	}
}

func TestPriorityString(t *testing.T) {
	testCases := []struct {
		priority Priority
		expected string
	}{
		{Low, "Low"},
		{Medium, "Medium"},
		{High, "High"},
	}

	for _, tc := range testCases {
		if tc.priority.String() != tc.expected {
			t.Errorf("Expected %s for priority %d, but got %s", tc.expected, tc.priority, tc.priority.String())
		}
	}
}

func TestTodosEmpty(t *testing.T) {
	todos := &Todos{}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print(todos)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	output := buf.String()
	if !strings.Contains(output, "No tasks. Your todo list is empty.") {
		t.Errorf("Expected 'No tasks' message for empty todo list, but got: %s", output)
	}
}

func TestVisualization(t *testing.T) {
	todos := &Todos{}
	todos.Add("Task 1", nil, High, nil)
	todos.Add("Task 2", nil, Medium, nil)
	todos.Add("Task 3", nil, Low, nil)
	todos.Add("Task 4", nil, Medium, nil)
	todos.Complete(0)

	priorityViz := VisualizeTasksByPriority(todos)
	if !strings.Contains(priorityViz, "High") || !strings.Contains(priorityViz, "Medium") || !strings.Contains(priorityViz, "Low") {
		t.Errorf("Priority visualization doesn't contain all priority levels")
	}

	progressViz := VisualizeOverallProgress(todos)
	if !strings.Contains(progressViz, "25.0%") {
		t.Errorf("Progress visualization doesn't show correct percentage")
	}
}
