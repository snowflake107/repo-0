package steps

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	projects2 "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/mcasperson/OctoterraWizard/internal/infrastructure"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/octoerrors"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"github.com/samber/lo"
	"net/url"
)

type StartProjectExportStep struct {
	BaseStep
	Wizard         wizard.Wizard
	exportProjects *widget.Button
	logs           *widget.Entry
	exportDone     bool
}

func (s StartProjectExportStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, previous, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(StartSpaceExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		moveNext := func(proceed bool) {
			if !proceed {
				return
			}

			s.Wizard.ShowWizardStep(FinishStep{
				Wizard:   s.Wizard,
				BaseStep: BaseStep{State: s.State}})
		}
		if !s.exportDone {
			dialog.NewConfirm(
				"Do you want to skip this step?",
				"You can run the runbooks manually from the Octopus UI.", moveNext, s.Wizard.Window).Show()
		} else {
			moveNext(true)
		}
	})
	linkUrl, _ := url.Parse(s.State.Server + "/app#/" + s.State.Space + "/tasks")
	link := widget.NewHyperlink("View the task list", linkUrl)
	link.Hide()
	s.logs = widget.NewEntry()
	s.logs.SetMinRowsVisible(20)
	s.logs.Disable()
	s.logs.Hide()
	s.logs.MultiLine = true
	s.exportDone = false

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The projects in the source space are now ready to begin exporting to the destination space.
		This involves serializing the project level resources (project, runbooks, variables, triggers etc) to a Terraform module and then applying the module to the destination space.
		First, this wizard publishes and runs the "__ 1. Serialize Project" runbook in each project to create the Terraform module.
		Then this wizard publishes and runs the "__ 2. Deploy Project" runbook in each project to apply the Terraform module to the destination space.
		Click the "Export Projects" button to execute these runbooks.
	`))
	result := widget.NewLabel("")
	infinite := widget.NewProgressBarInfinite()
	infinite.Hide()
	infinite.Start()
	s.exportProjects = widget.NewButton("Export Projects", func() {
		s.exportProjects.Disable()
		next.Disable()
		previous.Disable()
		infinite.Show()
		link.Hide()
		s.exportDone = true
		defer s.exportProjects.Enable()
		defer previous.Enable()
		defer next.Enable()
		defer infinite.Hide()

		result.SetText("🔵 Running the runbooks.")

		if err := s.Execute(func(message string) {
			result.SetText(message)
		}); err != nil {
			result.SetText(fmt.Sprintf("🔴 Failed to publish and run the runbooks. The failed tasks are shown below. You can review the task details in the Octopus console to find more information."))
			s.logs.SetText(err.Error())
			s.logs.Show()
			link.Show()
		} else {
			result.SetText("🟢 Runbooks ran successfully.")
			next.Enable()
			s.logs.Hide()
		}
	})
	middle := container.New(layout.NewVBoxLayout(), label1, s.exportProjects, infinite, result, link, s.logs)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s StartProjectExportStep) Execute(statusCallback func(message string)) error {
	myclient, err := octoclient.CreateClient(s.State)

	if err != nil {
		return err
	}

	projects, err := projects2.GetAll(myclient, myclient.GetSpaceID())

	if err != nil {
		return err
	}

	filteredProjects := lo.Filter(projects, func(project *projects2.Project, index int) bool {
		return project.Name != "Octoterra Space Management"
	})

	runAndTaskError := s.serializeProjects(filteredProjects, statusCallback)
	runAndTaskError = errors.Join(runAndTaskError, s.deployProjects(filteredProjects, statusCallback))

	return runAndTaskError
}

func (s StartProjectExportStep) serializeProjects(filteredProjects []*projects2.Project, statusCallback func(message string)) error {
	var runAndTaskError error = nil

	for _, project := range filteredProjects {
		if err := infrastructure.PublishRunbook(s.State, "__ 1. Serialize Project", project.Name); err != nil {
			return err
		}

		statusCallback("🔵 Published __ 1. Serialize Project runbook in project " + project.Name)
	}

	tasks := map[string]string{}
	for _, project := range filteredProjects {
		if taskId, err := infrastructure.RunRunbook(s.State, "__ 1. Serialize Project", project.Name); err != nil {

			var failedRunbookRun octoerrors.RunbookRunFailedError
			if errors.As(err, &failedRunbookRun) {
				runAndTaskError = errors.Join(runAndTaskError, failedRunbookRun)
			} else {
				return err
			}
		} else {
			tasks[project.Name] = taskId
		}
	}

	serializeIndex := 0
	for project, taskId := range tasks {
		if err := infrastructure.WaitForTask(s.State, taskId, func(message string) {
			statusCallback("🔵 __ 1. Serialize Project for project " + project + " is " + message + " (" + fmt.Sprint(serializeIndex) + "/" + fmt.Sprint(len(tasks)) + ")")
		}); err != nil {
			runAndTaskError = errors.Join(runAndTaskError, err)
		}
		serializeIndex++
	}

	return runAndTaskError
}

func (s StartProjectExportStep) deployProjects(filteredProjects []*projects2.Project, statusCallback func(message string)) error {
	var runAndTaskError error = nil

	for _, project := range filteredProjects {
		if err := infrastructure.PublishRunbook(s.State, "__ 2. Deploy Project", project.Name); err != nil {
			return err
		}
		statusCallback("🔵 Published __ 2. Deploy Space runbook in project " + project.Name)
	}

	applyTasks := map[string]string{}
	for _, project := range filteredProjects {
		if taskId, err := infrastructure.RunRunbook(s.State, "__ 2. Deploy Project", project.Name); err != nil {
			var failedRunbookRun octoerrors.RunbookRunFailedError
			if errors.As(err, &failedRunbookRun) {
				runAndTaskError = errors.Join(runAndTaskError, failedRunbookRun)
			} else {
				return err
			}
		} else {
			applyTasks[project.Name] = taskId
		}
	}

	applyIndex := 0
	for project, taskId := range applyTasks {
		if err := infrastructure.WaitForTask(s.State, taskId, func(message string) {
			statusCallback("🔵 __ 2. Deploy Project for project " + project + " is " + message + " (" + fmt.Sprint(applyIndex) + "/" + fmt.Sprint(len(applyTasks)) + ")")
		}); err != nil {
			runAndTaskError = errors.Join(runAndTaskError, err)
		}
		applyIndex++
	}

	return runAndTaskError
}
