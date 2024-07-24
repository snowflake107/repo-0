package octoerrors

import "fmt"

type TaskFailedError struct {
	TaskId string
}

func (e TaskFailedError) Error() string {
	return fmt.Sprintf("Task %s failed", e.TaskId)
}
