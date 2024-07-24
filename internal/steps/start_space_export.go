package steps

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/infrastructure"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type StartSpaceExportStep struct {
	BaseStep
	Wizard      wizard.Wizard
	exportSpace *widget.Button
}

func (s StartSpaceExportStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, previous, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(ProjectExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(StartProjectExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})
	next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The source space is now ready to begin exporting to the destination space.
		We start by serializing the space level resources (feeds, accounts, tenants, certificates, targets etc) using the runbooks in the "Octoterra Space Management" project.
		First, we run the "__ 1. Serialize Space" runbook to create the Terraform module.
		Then we run the "__ 2. Deploy Space" runbook to apply the Terraform module to the destination space.
		Click the "Export Space" button to execute these runbooks.
	`))
	result := widget.NewLabel("")
	infinite := widget.NewProgressBarInfinite()
	infinite.Hide()
	infinite.Start()
	s.exportSpace = widget.NewButton("Export Space", func() {
		s.exportSpace.Disable()
		next.Disable()
		previous.Disable()
		infinite.Show()
		defer s.exportSpace.Enable()
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
	middle := container.New(layout.NewVBoxLayout(), label1, s.exportSpace, infinite, result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s StartSpaceExportStep) Execute(statusCallback func(message string)) (executeError error) {
	if err := infrastructure.PublishRunbook(s.State, "__ 1. Serialize Space", "Octoterra Space Management"); err != nil {
		return err
	}

	statusCallback("🔵 Published __ 1. Serialize Space runbook")

	if taskId, err := infrastructure.RunRunbook(s.State, "__ 1. Serialize Space", "Octoterra Space Management"); err != nil {
		return err
	} else {
		if err := infrastructure.WaitForTask(s.State, taskId, func(message string) {
			statusCallback("🔵 __ 1. Serialize Space is " + message)
		}); err != nil {
			return err
		}
	}

	if err := infrastructure.PublishRunbook(s.State, "__ 2. Deploy Space", "Octoterra Space Management"); err != nil {
		return err
	}

	statusCallback("🔵 Published __ 2. Deploy Space runbook")

	if taskId, err := infrastructure.RunRunbook(s.State, "__ 2. Deploy Space", "Octoterra Space Management"); err != nil {
		return err
	} else {
		if err := infrastructure.WaitForTask(s.State, taskId, func(message string) {
			statusCallback("🔵 __ 2. Deploy Space is " + message + ". This runbook can take quite some time (many hours) for large spaces.")
		}); err != nil {
			return err
		}
	}

	return nil
}
