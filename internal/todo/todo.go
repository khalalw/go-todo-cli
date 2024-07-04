package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Todo struct {
	Task      string
	Completed bool
	DueDate   *time.Time `json:",omitempty"`
}

type Todos []Todo

func (t *Todos) Add(task string, dueDate *time.Time) {
	todo := Todo{Task: task, Completed: false, DueDate: dueDate}
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
		dueDate := ""
		if todo.DueDate != nil {
			dueDate = todo.DueDate.Format("2006-01-02")
		}
		fmt.Printf("%d. [%s] %s %s\n", i+1, status, todo.Task, dueDate)
	}
}
