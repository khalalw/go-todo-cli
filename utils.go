package main

import (
	"fmt"
)

func parseIndex(input string) int {
	var index int
	_, err := fmt.Sscanf(input, "%d", &index)
	if err != nil {
		return -1
	}
	return index - 1 // Convert to zero-based index
}

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  add <task>           - Add a new task")
	fmt.Println("  complete <task_number> - Mark a task as complete")
	fmt.Println("  delete <task_number> - Delete a task")
	fmt.Println("  list                 - List all tasks")
	fmt.Println("  help                 - Show this help message")
	fmt.Println("  exit                 - Exit the TODO CLI")
}
