package main

import (
	"fmt"
	"go-todo-cli/internal/shell"
	"go-todo-cli/pkg/todo"
	"os"
)

const fileToWrite = "todos.json"

func main() {
	todoList := &todo.Todos{}
	handleFileLoading(todoList, fileToWrite)

	// Start the interactive shell
	shell.StartShell(todoList)

	// Save changes before exiting
	if err := todoList.Save(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving todo list:", err)
		os.Exit(1)
	}
}

func handleFileLoading(todoList *todo.Todos, filename string) {
	if err := todoList.Load(filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found, creating a new todos.json file.")
			if err := todoList.Save(filename); err != nil {
				fmt.Fprintln(os.Stderr, "Error creating new todo file:", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error loading todo file:", err)
			os.Exit(1)
		}
	}
}
