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

func AddCommand(args []string, dueDate *time.Time, priority todo.Priority, todoList *todo.Todos, tags []string) {
	if len(args) < 1 {
		fmt.Println("Usage: add <task> [--tag tag1,tag2,...]")
		return
	}
	todoList.Add(strings.Join(args, " "), dueDate, priority, tags)
	fmt.Println("Task added.")
	saveTodoList(todoList)
}

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

func ListCommand(todoList *todo.Todos) {
	todo.Print(todoList)
}

func ClearTasksCommand(todoList *todo.Todos) {
	*todoList = todo.Todos{}
	fmt.Println("All tasks cleared.")
	saveTodoList(todoList)
}

func EditCommand(taskNumber int, todoList *todo.Todos) {
	index := taskNumber - 1
	if index < 0 || index >= len(*todoList) {
		fmt.Println("Invalid task number.")
		return
	}

	task := &(*todoList)[index]
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Current task: %s\nEnter new task description (or press Enter to keep current): ", task.Task)
	newTask, _ := reader.ReadString('\n')
	newTask = strings.TrimSpace(newTask)
	if newTask != "" {
		task.Task = newTask
	}

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

	fmt.Printf("Current tags: %v\nEnter new tags (comma-separated) or press Enter to keep current: ", task.Tags)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		task.Tags = strings.Split(input, ",")
	}

	fmt.Println("Task updated successfully.")
	saveTodoList(todoList)
}

func AddTagCommand(args []string, todoList *todo.Todos) {
	if len(args) != 2 {
		fmt.Println("Usage: add-tag <task_number> <tag>")
		return
	}
	index := parseIndex(args[0])
	if index >= 0 && index < len(*todoList) {
		task := &(*todoList)[index]
		newTag := strings.TrimSpace(args[1])
		if !contains(task.Tags, newTag) {
			task.Tags = append(task.Tags, newTag)
			fmt.Printf("Tag '%s' added to task %d.\n", newTag, index+1)
			saveTodoList(todoList)
		} else {
			fmt.Printf("Tag '%s' already exists for task %d.\n", newTag, index+1)
		}
	} else {
		fmt.Println("Invalid task number.")
	}
}

func RemoveTagCommand(args []string, todoList *todo.Todos) {
	if len(args) != 2 {
		fmt.Println("Usage: remove-tag <task_number> <tag>")
		return
	}
	index := parseIndex(args[0])
	if index >= 0 && index < len(*todoList) {
		task := &(*todoList)[index]
		tagToRemove := strings.TrimSpace(args[1])
		if removeString(&task.Tags, tagToRemove) {
			fmt.Printf("Tag '%s' removed from task %d.\n", tagToRemove, index+1)
			saveTodoList(todoList)
		} else {
			fmt.Printf("Tag '%s' not found for task %d.\n", tagToRemove, index+1)
		}
	} else {
		fmt.Println("Invalid task number.")
	}
}

func FilterByTagCommand(args []string, todoList *todo.Todos) {
	if len(args) != 1 {
		fmt.Println("Usage: filter-tag <tag>")
		return
	}
	tag := strings.TrimSpace(args[0])
	filteredList := todo.Todos{}
	for _, task := range *todoList {
		if contains(task.Tags, tag) {
			filteredList = append(filteredList, task)
		}
	}
	if len(filteredList) > 0 {
		todo.Print(&filteredList)
	} else {
		fmt.Printf("No tasks found with tag '%s'.\n", tag)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeString(slice *[]string, item string) bool {
	for i, s := range *slice {
		if s == item {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return true
		}
	}
	return false
}

func saveTodoList(todoList *todo.Todos) {
	if err := todoList.Save(FileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving go-todo-cli list:", err)
	}
}

func parseIndex(input string) int {
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index <= 0 {
		return -1 // Return -1 for invalid or non-positive inputs
	}
	return index - 1 // Convert to zero-based index
}
