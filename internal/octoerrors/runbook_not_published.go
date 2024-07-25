package octoerrors

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
)

type RunbookNotPublishedError struct {
	Runbook *runbooks.Runbook
	Project *projects.Project
}

func (e RunbookNotPublishedError) Error() string {
	return fmt.Sprintf("Runbook %s %s is not published", e.Project.Name, e.Runbook.Name)
}
