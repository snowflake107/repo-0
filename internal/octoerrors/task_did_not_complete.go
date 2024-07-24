package octoerrors

import "fmt"

type TaskDidNotCompleteError struct {
	TaskId string
}

func (e TaskDidNotCompleteError) Error() string {
	return fmt.Sprintf("Task %s did not complete", e.TaskId)
}
