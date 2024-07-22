package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/query"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type StepTemplateStep struct {
	BaseStep
	Wizard wizard.Wizard
	result *widget.Label
	logs   *widget.Entry
}

func (s StepTemplateStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(SpreadVariablesStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(SpaceExportStep{Wizard: s.Wizard, BaseStep: BaseStep{State: s.State}})
	})
	next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The runbooks created by this wizard require a number of step templates to be installed from the community step template library.
	`))
	s.result = widget.NewLabel("")
	s.logs = widget.NewEntry()
	s.logs.Disable()
	s.logs.MultiLine = true
	s.logs.SetMinRowsVisible(20)
	s.logs.Hide()

	installSteps := widget.NewButton("Install Step Templates", func() {
		s.logs.Hide()
		s.result.SetText("🔵 Installing step templates.")
		myclient, err := octoclient.CreateClient(s.State)

		if err != nil {
			s.logs.SetText("🔴 Failed to create the client:\n" + err.Error())
			return
		}

		// Octopus - Serialize Space to Terraform
		if err, message := query.InstallStepTemplate(myclient, s.State, "https://library.octopus.com/step-templates/e03c56a4-f660-48f6-9d09-df07e1ac90bd"); err != nil {
			s.result.SetText(message)
			s.logs.Show()
			s.logs.SetText(err.Error())
			return
		}

		// Octopus - Populate Octoterra Space (S3 Backend)
		if err, message := query.InstallStepTemplate(myclient, s.State, "https://library.octopus.com/step-templates/14d51af4-1c3d-4d41-9044-4304111d0cd8"); err != nil {
			s.result.SetText(message)
			s.logs.Show()
			s.logs.SetText(err.Error())
			return
		}

		next.Enable()
		s.result.SetText("🟢 Step templates installed.")
	})
	middle := container.New(layout.NewVBoxLayout(), label1, installSteps, s.result, s.logs)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
