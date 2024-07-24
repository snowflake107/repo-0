package steps

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	projects2 "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/mcasperson/OctoterraWizard/internal/infrastructure"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"github.com/samber/lo"
)

type StartProjectExportStep struct {
	BaseStep
	Wizard         wizard.Wizard
	exportProjects *widget.Button
}

func (s StartProjectExportStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, previous, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(StartSpaceExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(FinishStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})
	next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The projects in the source space are now ready to begin exporting to the destination space.
		We start by serializing the project level resources (project, runbooks, variables, triggers etc) using two runbooks added to each project.
		First, we run the "__ 1. Serialize Project" runbook to create the Terraform module.
		Then we run the "__ 2. Deploy Project" runbook to apply the Terraform module to the destination space.
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
		defer s.exportProjects.Enable()
		defer previous.Enable()
		defer infinite.Hide()

		result.SetText("🔵 Running the runbooks.")

		if err := s.Execute(func(message string) {
			result.SetText(message)
		}); err != nil {
			result.SetText(fmt.Sprintf("🔴 Failed to publish and run the runbooks: %s", err))
		} else {
			result.SetText("🟢 Runbooks ran successfully.")
			next.Enable()
		}
	})
	middle := container.New(layout.NewVBoxLayout(), label1, s.exportProjects, infinite, result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s StartProjectExportStep) Execute(statusCallback func(message string)) (executeError error) {
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

	for _, project := range filteredProjects {
		if err := infrastructure.PublishRunbook(s.State, "__ 1. Serialize Project", project.Name); err != nil {
			return err
		}

		statusCallback("🔵 Published __ 1. Serialize Project runbook in project " + project.Name)
	}

	tasks := map[string]string{}

	for _, project := range filteredProjects {
		if taskId, err := infrastructure.RunRunbook(s.State, "__ 1. Serialize Project", project.Name); err != nil {
			return err
		} else {
			tasks[project.Name] = taskId
		}
	}

	serializeIndex := 0
	for project, taskId := range tasks {
		if err := infrastructure.WaitForTask(s.State, taskId, func(message string) {
			statusCallback("🔵 __ 1. Serialize Project for project " + project + " is " + message + " (" + fmt.Sprint(serializeIndex) + "/" + fmt.Sprint(len(tasks)) + ")")
		}); err != nil {
			return err
		}
		serializeIndex++
	}

	for _, project := range filteredProjects {
		if err := infrastructure.PublishRunbook(s.State, "__ 2. Deploy Project", project.Name); err != nil {
			return err
		}
		statusCallback("🔵 Published __ 2. Deploy Space runbook in project " + project.Name)
	}

	applyTasks := map[string]string{}
	for _, project := range filteredProjects {
		if taskId, err := infrastructure.RunRunbook(s.State, "__ 2. Deploy Space", "Octoterra Space Management"); err != nil {
			return err
		} else {
			tasks[project.Name] = taskId
		}
	}

	applyIndex := 0
	for project, taskId := range applyTasks {
		if err := infrastructure.WaitForTask(s.State, taskId, func(message string) {
			statusCallback("🔵 __ 2. Deploy Space for project " + project + " is " + message + " (" + fmt.Sprint(applyIndex) + "/" + fmt.Sprint(len(applyTasks)) + ")")
		}); err != nil {
			return err
		}
		applyIndex++
	}

	return nil
}
