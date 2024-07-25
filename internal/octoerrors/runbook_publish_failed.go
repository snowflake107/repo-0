package octoerrors

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
)

type RunbookPublishFailedError struct {
	Runbook  *runbooks.Runbook
	Project  *projects.Project
	Response string
}

func (e RunbookPublishFailedError) Error() string {
	return fmt.Sprintf("Runbook %s %s failed to publish: %s", e.Project.Name, e.Runbook.Name, e.Response)
}
