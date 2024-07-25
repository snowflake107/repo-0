package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tasks"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/octoerrors"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/samber/lo"
	"io"
	"net/http"
	"strings"
	"time"
)

func WaitForTask(state state.State, taskId string, statusCallback func(message string)) error {
	myclient, err := octoclient.CreateClient(state)

	if err != nil {
		return err
	}

	// wait up to 2 hours for the task to complete
	for i := 0; i < 7200; i++ {
		mytasks, err := myclient.Tasks.Get(tasks.TasksQuery{
			Environment:             "",
			HasPendingInterruptions: false,
			HasWarningsOrErrors:     false,
			IDs:                     []string{taskId},
			IncludeSystem:           false,
			IsActive:                false,
			IsRunning:               false,
			Name:                    "",
			Node:                    "",
			PartialName:             "",
			Project:                 "",
			Runbook:                 "",
			Skip:                    0,
			Spaces:                  nil,
			States:                  nil,
			Take:                    1,
			Tenant:                  "",
		})

		if err != nil {
			return err
		}

		if len(mytasks.Items) == 0 {
			return octoerrors.TaskNotFound{TaskId: taskId}
		}

		if mytasks.Items[0].IsCompleted != nil && *mytasks.Items[0].IsCompleted {
			if mytasks.Items[0].State != "Success" {
				return octoerrors.TaskFailedError{TaskId: taskId}
			}
			statusCallback(mytasks.Items[0].State)
			return nil
		} else {
			statusCallback(mytasks.Items[0].State)
			time.Sleep(10 * time.Second)
		}
	}

	return octoerrors.TaskDidNotCompleteError{TaskId: taskId}
}

func RunRunbook(state state.State, runbookName string, projectName string) (string, error) {
	myclient, err := octoclient.CreateClient(state)

	if err != nil {
		return "", err
	}

	environment, err := environments.GetAll(myclient, myclient.GetSpaceID())

	if err != nil {
		return "", err
	}

	project, err := projects.GetByName(myclient, myclient.GetSpaceID(), projectName)

	if err != nil {
		return "", err
	}

	runbook, err := runbooks.GetByName(myclient, myclient.GetSpaceID(), project.GetID(), runbookName)

	if err != nil {
		return "", err
	}

	url := state.Server + runbook.GetLinks()["RunbookRunPreview"]
	url = strings.ReplaceAll(url, "{environment}", environment[0].GetID())
	url = strings.ReplaceAll(url, "{?includeDisabledSteps}", "")

	runbookRunPreviewRequest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	runbookRunPreviewResponse, err := myclient.HttpSession().DoRawRequest(runbookRunPreviewRequest)

	if err != nil {
		return "", err
	}

	runbookRunPreviewRaw, err := io.ReadAll(runbookRunPreviewResponse.Body)

	if err != nil {
		return "", err
	}

	runbookRunPreview := map[string]any{}
	err = json.Unmarshal(runbookRunPreviewRaw, &runbookRunPreview)

	if err != nil {
		return "", err
	}

	runbookFormNames := lo.Map(runbookRunPreview["Form"].(map[string]any)["Elements"].([]any), func(value any, index int) any {
		return value.(map[string]any)["Name"]
	})

	runbookFormValues := map[string]string{}

	for _, name := range runbookFormNames {
		runbookFormValues[name.(string)] = "dummy"
	}

	runbookBody := map[string]any{
		"RunbookId":                runbook.GetID(),
		"RunbookSnapShotId":        runbook.PublishedRunbookSnapshotID,
		"FrozenRunbookProcessId":   nil,
		"EnvironmentId":            environment[0].GetID(),
		"TenantId":                 nil,
		"SkipActions":              []string{},
		"QueueTime":                nil,
		"QueueTimeExpiry":          nil,
		"FormValues":               runbookFormValues,
		"ForcePackageDownload":     false,
		"ForcePackageRedeployment": true,
		"UseGuidedFailure":         false,
		"SpecificMachineIds":       []string{},
		"ExcludedMachineIds":       []string{},
	}

	runbookBodyJson, err := json.Marshal(runbookBody)

	if err != nil {
		return "", err
	}

	url = state.Server + "/api/" + state.Space + "/runbookRuns"
	runbookRunRequest, err := http.NewRequest("POST", url, bytes.NewReader(runbookBodyJson))

	if err != nil {
		return "", err
	}

	runbookRunResponse, err := myclient.HttpSession().DoRawRequest(runbookRunRequest)

	if err != nil {
		return "", err
	}

	runbookRunRaw, err := io.ReadAll(runbookRunResponse.Body)

	if err != nil {
		return "", err
	}

	runbookRun := map[string]any{}
	err = json.Unmarshal(runbookRunRaw, &runbookRun)

	if err != nil {
		return "", err
	}

	if _, ok := runbookRun["TaskId"]; !ok {
		return "", octoerrors.RunbookRunFailedError{Runbook: runbook, Project: project, Response: string(runbookRunRaw)}
	}

	return runbookRun["TaskId"].(string), nil

}

func PublishRunbook(state state.State, runbookName string, projectName string) error {
	myclient, err := octoclient.CreateClient(state)

	if err != nil {
		return err
	}

	project, err := projects.GetByName(myclient, myclient.GetSpaceID(), projectName)

	if err != nil {
		return err
	}

	runbook, err := runbooks.GetByName(myclient, myclient.GetSpaceID(), project.GetID(), runbookName)

	if err != nil {
		return err
	}

	url := state.Server + runbook.GetLinks()["RunbookSnapshotTemplate"]
	runbookSnapshotTemplateRequest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	runbookSnapshotTemplateResponse, err := myclient.HttpSession().DoRawRequest(runbookSnapshotTemplateRequest)

	if err != nil {
		return err
	}

	runbookSnapshotTemplateRaw, err := io.ReadAll(runbookSnapshotTemplateResponse.Body)

	if err != nil {
		return err
	}

	runbookSnapshotTemplate := map[string]any{}
	err = json.Unmarshal(runbookSnapshotTemplateRaw, &runbookSnapshotTemplate)

	if err != nil {
		return err
	}

	snapshot := map[string]any{
		"ProjectId": project.GetID(),
		"RunbookId": runbook.GetID(),
		"Name":      runbookSnapshotTemplate["NextNameIncrement"],
	}

	var packageErrors error = nil
	snapshot["SelectedPackages"] = lo.Map(runbookSnapshotTemplate["Packages"].([]any), func(pkg any, index int) any {
		snapshotPackage := pkg.(map[string]any)
		versions, err := feeds.SearchPackageVersions(myclient, myclient.GetSpaceID(), snapshotPackage["FeedId"].(string), snapshotPackage["PackageId"].(string), "", 1)

		if err != nil {
			packageErrors = errors.Join(packageErrors, err)
			return nil
		}

		return map[string]any{
			"ActionName":           snapshotPackage["ActionName"],
			"Version":              versions.Items[0].Version,
			"PackageReferenceName": snapshotPackage["PackageReferenceName"],
		}
	})

	if packageErrors != nil {
		return packageErrors
	}

	snapshotJson, err := json.Marshal(snapshot)

	if err != nil {
		return err
	}

	url = state.Server + "/api/" + state.Space + "/runbookSnapshots?publish=true"
	runbookSnapshotRequest, err := http.NewRequest("POST", url+"?publish=true", bytes.NewBuffer(snapshotJson))

	if err != nil {
		return err
	}

	runbookSnapshotResponse, err := myclient.HttpSession().DoRawRequest(runbookSnapshotRequest)

	if err != nil {
		return err
	}

	runbookSnapshotResponseRaw, err := io.ReadAll(runbookSnapshotResponse.Body)

	if err != nil {
		return err
	}

	runbookSnapshot := map[string]any{}
	err = json.Unmarshal(runbookSnapshotResponseRaw, &runbookSnapshot)

	if err != nil {
		return err
	}

	fmt.Println(runbookSnapshot)

	return nil
}
