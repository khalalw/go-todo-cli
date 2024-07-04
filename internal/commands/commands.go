package commands

import (
	"bufio"
	"fmt"
	"go-todo-cli/internal/todo"
	"os"
	"strconv"
	"strings"
	"time"
)

var FileToWrite = "todos.json"

// AddCommand adds a new task to the TODO list
func AddCommand(args []string, dueDate *time.Time, priority todo.Priority, todoList *todo.Todos) {
	if len(args) < 1 {
		fmt.Println("Usage: add <task>")
		return
	}
	todoList.Add(strings.Join(args, " "), dueDate, priority)
	fmt.Println("Task added.")
	saveTodoList(todoList)
}

// CompleteCommand marks a task as complete
func CompleteCommand(args []string, todoList *todo.Todos) {
	if len(args) != 1 {
		fmt.Println("Usage: complete <task_number>")
		return
	}
	index := parseIndex(args[0])
	if index >= 0 && index < len(*todoList) {
		if err := todoList.Complete(index); err == nil {
			fmt.Println("Task marked as complete.")
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Invalid task number.")
	}
	saveTodoList(todoList)
}

// DeleteCommand deletes a task from the TODO list
func DeleteCommand(args []string, todoList *todo.Todos) {
	if len(args) != 1 {
		fmt.Println("Usage: delete <task_number>")
		return
	}
	index := parseIndex(args[0])
	if index >= 0 && index < len(*todoList) {
		if err := todoList.Delete(index); err == nil {
			fmt.Println("Task deleted.")
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Invalid task number.")
	}
	saveTodoList(todoList)
}

// ListCommand lists all tasks
func ListCommand(todoList *todo.Todos) {
	todo.Print(todoList)
}

// ClearTasksCommand clears all tasks from the TODO list
func ClearTasksCommand(todoList *todo.Todos) {
	*todoList = todo.Todos{}
	fmt.Println("All tasks cleared.")
	saveTodoList(todoList)
}

// EditCommand allows editing an existing task
func EditCommand(taskNumber int, todoList *todo.Todos) {
	index := taskNumber - 1
	if index < 0 || index >= len(*todoList) {
		fmt.Println("Invalid task number.")
		return
	}

	task := &(*todoList)[index]
	reader := bufio.NewReader(os.Stdin)

	// Edit task description
	fmt.Printf("Current task: %s\nEnter new task description (or press Enter to keep current): ", task.Task)
	newTask, _ := reader.ReadString('\n')
	newTask = strings.TrimSpace(newTask)
	if newTask != "" {
		task.Task = newTask
	}

	// Edit due date
	var newDueDate time.Time
	var err error
	for {
		fmt.Printf("Current due date: %v\nEnter new due date (YYYY-MM-DD) or press Enter to keep current: ", task.DueDate)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			break
		}
		newDueDate, err = time.Parse("2006-01-02", input)
		if err == nil {
			task.DueDate = &newDueDate
			break
		}
		fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
	}

	// Edit priority
	for {
		fmt.Printf("Current priority: %v\nEnter new priority (low/medium/high) or press Enter to keep current: ", task.Priority)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			break
		}
		newPriority, err := todo.ParsePriority(strings.ToLower(input))
		if err == nil {
			task.Priority = newPriority
			break
		}
		fmt.Println("Invalid priority. Please use low, medium, or high.")
	}

	fmt.Println("Task updated successfully.")
	saveTodoList(todoList)
}

// saveTodoList saves the TODO list to a file
func saveTodoList(todoList *todo.Todos) {
	if err := todoList.Save(FileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving go-todo-cli list:", err)
	}
}

// parseIndex parses a string index to an integer
func parseIndex(input string) int {
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index <= 0 {
		return -1 // Return -1 for invalid or non-positive inputs
	}
	return index - 1 // Convert to zero-based index
}
