package main

import (
	"fmt"
	"go-todo-cli/todo"
	"os"
	"strings"
)

func executeCommand(cmd string, args []string, todoList *todo.Todos) {
	switch cmd {
	case "add":
		addCommand(args, todoList)
	case "complete":
		completeCommand(args, todoList)
	case "delete":
		deleteCommand(args, todoList)
	case "list":
		listCommand(todoList)
	case "exit":
		exitCommand(todoList)
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command. Type 'help' for a list of commands.")
	}
}

func addCommand(args []string, todoList *todo.Todos) {
	if len(args) < 1 {
		fmt.Println("Usage: add <task>")
		return
	}
	todoList.Add(strings.Join(args, " "))
	fmt.Println("Task added.")
	saveTodoList(todoList)
}

func completeCommand(args []string, todoList *todo.Todos) {
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

func deleteCommand(args []string, todoList *todo.Todos) {
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

func listCommand(todoList *todo.Todos) {
	todo.Print(todoList)
}

func exitCommand(todoList *todo.Todos) {
	fmt.Println("Exiting the TODO CLI.")
	saveTodoList(todoList)
	os.Exit(0)
}

func saveTodoList(todoList *todo.Todos) {
	if err := todoList.Save(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving todo list:", err)
	}
}
