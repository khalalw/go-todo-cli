package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
	sort.Slice(*todos, func(i, j int) bool {
		if (*todos)[i].Priority == (*todos)[j].Priority {
			return (*todos)[i].Task < (*todos)[j].Task
		}
		return (*todos)[i].Priority > (*todos)[j].Priority
	})

	for i, todo := range *todos {
		status := " "
		if todo.Completed {
			status = "x"
		}
		dueDate := ""
		if todo.DueDate != nil {
			dueDate = todo.DueDate.Format("2006-01-02")
		}
		fmt.Printf("%d. [%s] %s %s (Priority: %s)\n", i+1, status, todo.Task, dueDate, todo.Priority)
	}
}
