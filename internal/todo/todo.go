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
	switch strings.ToLower(s) {
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
	Tags      []string   `json:",omitempty"`
}

type Todos []Todo

func (t *Todos) Add(task string, dueDate *time.Time, priority Priority, tags []string) {
	todo := Todo{Task: task, Completed: false, DueDate: dueDate, Priority: priority, Tags: tags}
	*t = append(*t, todo)
}

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
	data, err := json.MarshalIndent(t, "", "  ")
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

	maxTaskLength := 0
	for _, todo := range *todos {
		if len(todo.Task) > maxTaskLength {
			maxTaskLength = len(todo.Task)
		}
	}

	format := "| %-3d | %-*s | %-10s | %-8s | %-6s | %-20s |\n"
	divider := strings.Repeat("-", maxTaskLength+60)

	fmt.Println(divider)
	fmt.Printf(format, 0, maxTaskLength, "Task", "Due Date", "Priority", "Status", "Tags")
	fmt.Println(divider)

	for i, todo := range *todos {
		status := "Pending"
		if todo.Completed {
			status = "Done"
		}

		dueDate := "N/A"
		if todo.DueDate != nil {
			dueDate = todo.DueDate.Format("2006-01-02")
		}

		tags := "None"
		if len(todo.Tags) > 0 {
			tags = strings.Join(todo.Tags, ", ")
		}

		fmt.Printf(format, i+1, maxTaskLength, todo.Task, dueDate, todo.Priority, status, tags)
	}

	fmt.Println(divider)
}
