package octoerrors

import "fmt"

type TaskNotFound struct {
	TaskId string
}

func (e TaskNotFound) Error() string {
	return fmt.Sprintf("Task %s was not found", e.TaskId)
}
