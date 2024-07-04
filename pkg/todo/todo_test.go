package todo

import (
	"bytes"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	todoList := &Todos{}
	task := "Test task"
	todoList.Add(task)

	if len(*todoList) != 1 {
		t.Fatalf("Expected length 1, got %d", len(*todoList))
	}

	if (*todoList)[0].Task != task {
		t.Errorf("Expected task %q, got %q", task, (*todoList)[0].Task)
	}

	if (*todoList)[0].Completed {
		t.Error("Expected task to be not completed")
	}
}

func TestComplete(t *testing.T) {
	todoList := &Todos{
		{Task: "Task 1", Completed: false},
		{Task: "Task 2", Completed: false},
	}
	err := todoList.Complete(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !(*todoList)[1].Completed {
		t.Error("Expected task to be completed")
	}
}

func TestCompleteInvalidIndex(t *testing.T) {
	todoList := &Todos{
		{Task: "Task 1", Completed: false},
	}
	err := todoList.Complete(2)
	if err == nil {
		t.Error("Expected error for invalid index, got nil")
	}
}

func TestDelete(t *testing.T) {
	todoList := &Todos{
		{Task: "Task 1", Completed: false},
		{Task: "Task 2", Completed: false},
	}
	err := todoList.Delete(0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(*todoList) != 1 {
		t.Fatalf("Expected length 1, got %d", len(*todoList))
	}

	if (*todoList)[0].Task != "Task 2" {
		t.Errorf("Expected task %q, got %q", "Task 2", (*todoList)[0].Task)
	}
}

func TestDeleteInvalidIndex(t *testing.T) {
	todoList := &Todos{
		{Task: "Task 1", Completed: false},
	}
	err := todoList.Delete(2)
	if err == nil {
		t.Error("Expected error for invalid index, got nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	todoList := &Todos{
		{Task: "Task 1", Completed: false},
		{Task: "Task 2", Completed: true},
	}
	filename := "test_todos.json"
	defer os.Remove(filename)

	err := todoList.Save(filename)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	loadedTodoList := &Todos{}
	err = loadedTodoList.Load(filename)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(*loadedTodoList) != 2 {
		t.Fatalf("Expected length 2, got %d", len(*loadedTodoList))
	}

	if (*loadedTodoList)[0].Task != "Task 1" || (*loadedTodoList)[1].Task != "Task 2" {
		t.Errorf("Tasks do not match after loading")
	}

	if (*loadedTodoList)[0].Completed || !(*loadedTodoList)[1].Completed {
		t.Error("Completion status does not match after loading")
	}
}

func TestPrint(t *testing.T) {
	todoList := &Todos{
		{Task: "Task B", Completed: false},
		{Task: "Task A", Completed: true},
	}
	expectedOutput := "1. [x] Task A\n2. [ ] Task B\n"

	actualOutput := captureOutput(func() {
		Print(todoList)
	})

	if actualOutput != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, actualOutput)
	}
}

// Helper function to capture output
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
