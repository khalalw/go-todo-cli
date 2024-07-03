package todo

import "fmt"

type Todo struct {
	Task      string
	Completed bool
}

type Todos []Todo

// Add updates our todos list with a new item, defaulted to false
func (t *Todos) Add(task string) {
	todo := Todo{Task: task, Completed: false}
	*t = append(*t, todo)
}

// Complete mark the todo complete at a given index
func (t *Todos) Complete(index int) error {
	// check to see if index is in range
	if index < 0 || index >= len(*t) {
		return fmt.Errorf("index out of range")
	}

	*t = append((*t)[:index], (*t)[index+1:]...)
	return nil
}
