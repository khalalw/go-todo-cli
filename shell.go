package main

import (
	"bufio"
	"fmt"
	"go-todo-cli/todo"
	"os"
	"strings"
)

func startShell(todoList *todo.Todos) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the TODO CLI. Type 'help' for a list of commands.")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		args := strings.Split(input, " ")
		cmd := args[0]

		executeCommand(cmd, args[1:], todoList)
	}
}
