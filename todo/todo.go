package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

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

	(*t)[index].Completed = true
	return nil
}

// Delete removes a item at a gien index
func (t *Todos) Delete(index int) error {
	if index < 0 || index >= len(*t) {
		return fmt.Errorf("index out of range")
	}

	*t = append((*t)[:index], (*t)[index+1:]...)
	return nil
}

// Save method updates the todos
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
		return (*todos)[i].Task < (*todos)[j].Task
	})

	for i, todo := range *todos {
		status := " "
		if todo.Completed {
			status = "x"
		}

		fmt.Printf("%d. [%s] %s\n", i+1, status, todo.Task)
	}
}
