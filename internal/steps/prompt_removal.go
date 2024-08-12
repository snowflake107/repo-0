package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type PromptRemovalStep struct {
	BaseStep
	Wizard wizard.Wizard
}

func (s PromptRemovalStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, _ := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(StepTemplateStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(SpaceExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})

	heading := widget.NewLabel("Prompt for Deletion")
	heading.TextStyle = fyne.TextStyle{Bold: true}

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The actions taken in the subsequent steps may involve deleting and recreating resources in your Octopus Deploy instance.
		You can manually approve each deletion, or you can allow the wizard to delete resources without being prompted,
	`))

	radio := widget.NewRadioGroup([]string{"Prompt for each deletion", "Automatically delete"}, func(value string) {
		s.State.PromptForDelete = value == "Prompt for each deletion"
	})

	if s.State.PromptForDelete {
		radio.SetSelected("Prompt for each deletion")
	} else {
		radio.SetSelected("Automatically delete")

	}

	middle := container.New(layout.NewVBoxLayout(), heading, label1, radio)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
