package todo

import (
	"fmt"
	"strings"
)

const (
	barChar     = "█"
	emptyChar   = "░"
	maxBarWidth = 20
)

func VisualizeTasksByPriority(todos *Todos) string {
	priorities := make(map[Priority]int)
	for _, todo := range *todos {
		priorities[todo.Priority]++
	}

	maxCount := 0
	for _, count := range priorities {
		if count > maxCount {
			maxCount = count
		}
	}

	var result strings.Builder
	result.WriteString("Task Distribution by Priority:\n\n")

	for _, priority := range []Priority{Low, Medium, High} {
		count := priorities[priority]
		barWidth := int(float64(count) / float64(maxCount) * maxBarWidth)
		bar := strings.Repeat(barChar, barWidth) + strings.Repeat(emptyChar, maxBarWidth-barWidth)
		result.WriteString(fmt.Sprintf("%-6s |%s| %d\n", priority, bar, count))
	}

	return result.String()
}

func VisualizeOverallProgress(todos *Todos) string {
	total := len(*todos)
	if total == 0 {
		return "No tasks available."
	}

	completed := 0
	for _, todo := range *todos {
		if todo.Completed {
			completed++
		}
	}

	percentage := float64(completed) / float64(total) * 100
	barWidth := int(percentage / 5) // 20 characters for 100%
	bar := strings.Repeat(barChar, barWidth) + strings.Repeat(emptyChar, 20-barWidth)

	return fmt.Sprintf("Overall Progress:\n\n[%s] %.1f%% (%d/%d tasks completed)", bar, percentage, completed, total)
}
