package todo

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	todos := &Todos{}
	todos.Add("Test task")
	if len(*todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(*todos))
	}
	if (*todos)[0].Task != "Test task" {
		t.Errorf("Expected 'Test task', got '%s'", (*todos)[0].Task)
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
