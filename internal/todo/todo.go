package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

func (p Priority) String() string {
	return [...]string{"Low", "Medium", "High"}[p]
}

func ParsePriority(s string) (Priority, error) {
	switch s {
	case "low":
		return Low, nil
	case "medium":
		return Medium, nil
	case "high":
		return High, nil
	default:
		return Low, fmt.Errorf("invalid priority: %s", s)
	}
}

type Todo struct {
	Task      string
	Completed bool
	DueDate   *time.Time `json:",omitempty"`
	Priority  Priority   `json:",omitempty"`
}

type Todos []Todo

func (t *Todos) Add(task string, dueDate *time.Time, priority Priority) {
	todo := Todo{Task: task, Completed: false, DueDate: dueDate, Priority: priority}
	*t = append(*t, todo)
}

// Complete marks a task as complete
func (t *Todos) Complete(index int) error {
	if index < 0 || index >= len(*t) {
		return fmt.Errorf("index out of range")
	}
	(*t)[index].Completed = true
	return nil
}

func (t *Todos) Delete(index int) error {
	if index < 0 || index >= len(*t) {
		return fmt.Errorf("index out of range")
	}
	*t = append((*t)[:index], (*t)[index+1:]...)
	return nil
}

func (t *Todos) Save(filename string) error {
	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, t)
}

func Print(todos *Todos) {
	if len(*todos) == 0 {
		fmt.Println("No tasks. Your todo list is empty.")
		return
	}

	// Find the longest task for proper spacing
	maxTaskLength := 0
	for _, todo := range *todos {
		if len(todo.Task) > maxTaskLength {
			maxTaskLength = len(todo.Task)
		}
	}

	// Print header
	fmt.Println(strings.Repeat("-", maxTaskLength+50))
	fmt.Printf("| %-3s | %-*s | %-10s | %-8s | %-6s |\n", "No.", maxTaskLength, "Task", "Due Date", "Priority", "Status")
	fmt.Println(strings.Repeat("-", maxTaskLength+50))

	// Print tasks
	for i, todo := range *todos {
		status := "Pending"
		if todo.Completed {
			status = "Done"
		}

		dueDate := "N/A"
		if todo.DueDate != nil {
			dueDate = todo.DueDate.Format("2006-01-02")
		}

		fmt.Printf("| %-3d | %-*s | %-10s | %-8s | %-6s |\n",
			i+1, maxTaskLength, todo.Task, dueDate, todo.Priority, status)
	}

	fmt.Println(strings.Repeat("-", maxTaskLength+50))
}
