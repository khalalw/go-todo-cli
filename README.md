# Go TODO CLI

A simple command-line interface (CLI) for managing a TODO list written in Go. This CLI allows you to add tasks, mark tasks as complete, delete tasks, and list all tasks, with the tasks being saved to a JSON file.

## Features

- Add a new task
- Mark a task as complete
- Delete a task
- List all tasks

## Requirements

- Go 1.16 or later

## Installation

1. Clone the repository:
    ```shell
    git clone https://github.com/yourusername/go-todo-cli.git
    ```
2. Change into the project directory:
    ```shell
    cd go-todo-cli
    ```

## Usage

### Adding a Task

To add a new task, use the `-add` flag followed by the task description:

```shell
go run main.go -add "Write Go documentation"
````

### Marking a Task as Complete
To mark a task as complete, use the `-complete` flag followed by the task number:
```shell
go run main.go -complete 1
```

### Deleting a Task
To delete a task, use the `-delete` flag followed by the task number:
```shell
go run main.go -delete 1
```

### Listing All Tasks
```shell
go run main.go -list
```

## Example
Here's an example sequence of commands:

1. Add tasks:
   ```shell
   go run main.go -add "Buy groceries"
   go run main.go -add "Read a book"
   ```
   
2. List tasks:
   ```shell
   go run main.go -list
   ```
   Output:
   ```
   1. [ ] Buy groceries
   2. [ ] Read a book
   ```
   
3. Mark a task as complete:
   ```sh
   go run main.go -complete 1
   ```

4. List tasks again:
   ```shell
   go run main.go -list
   ```
   Output:
   ```
   1. [x] Buy groceries
   2. [ ] Read a book
   ```
   
5. Delete a task:
   ```shell
   go run main.go -delete 1
   ```
   Output:
   ```
   1. [ ] Read a book
   ```

## Author
Khalal Walker