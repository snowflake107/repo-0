package octoerrors

import (
	"fmt"
	"strings"
)

type FailedTasksError struct {
	TaskId []string
}

func (e FailedTasksError) Error() string {
	return fmt.Sprintf("Tasks %s failed", strings.Join(e.TaskId, ","))
}
