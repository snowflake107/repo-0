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
	logs        *widget.Entry
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

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The the source space are now ready to begin exporting to the destination space.
		This involves serializing the space level resources (feeds, accounts, targets, tenants etc) to a Terraform module and then applying the module to the destination space.
		First, this wizard publishes and runs the "__ 1. Serialize Space" runbook in the "Octoterra Space Management" project to create the Terraform module.
		Then this wizard publishes and runs the "__ 2. Deploy Space" runbook in the "Octoterra Space Management" project to apply the Terraform module to the destination space.
		Click the "Export Space" button to execute these runbooks.
	`))
	result := widget.NewLabel("")
	infinite := widget.NewProgressBarInfinite()
	infinite.Hide()
	infinite.Start()
	s.logs = widget.NewEntry()
	s.logs.SetMinRowsVisible(20)
	s.logs.Disable()
	s.logs.Hide()
	s.logs.MultiLine = true

	s.exportSpace = widget.NewButton("Export Space", func() {
		s.exportSpace.Disable()
		previous.Disable()
		next.Disable()
		infinite.Show()
		s.logs.Hide()

		result.SetText("🔵 Running the runbooks.")

		if err := s.Execute(func(message string) {
			result.SetText(message)
		}); err != nil {
			result.SetText(fmt.Sprintf("🔴 Failed to publish and run the runbooks"))
			s.logs.Show()
			s.logs.SetText(err.Error())
			next.Enable()
			previous.Enable()
			infinite.Hide()
			s.exportSpace.Enable()
		} else {
			result.SetText("🟢 Runbooks ran successfully.")
			next.Enable()
			previous.Enable()
			s.logs.Hide()
			infinite.Hide()
			s.exportSpace.Enable()
		}
	})
	middle := container.New(layout.NewVBoxLayout(), label1, s.exportSpace, infinite, result, s.logs)

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
