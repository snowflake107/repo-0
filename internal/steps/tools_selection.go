package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type ToolsSelectionStep struct {
	BaseStep
	Wizard wizard.Wizard
}

func (s ToolsSelectionStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, _ := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(OctopusDestinationDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(BackendSelectionStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})

	heading := widget.NewLabel("Tools Selection")
	heading.TextStyle = fyne.TextStyle{Bold: true}

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The source server must either use Docker and container images to expose tools like Python and Terraform,
		or have Terraform and Python installed locally.
		If your source server has Docker installed, select the "Container Images" option.
		Otherwise, ensure that Terraform and Python are installed locally and select the "Local Tools" option.
	`))

	radio := widget.NewRadioGroup([]string{"Container Images", "Local Tools"}, func(value string) {
		s.State.UseContainerImages = value == "Container Images"
	})

	if s.State.UseContainerImages {
		radio.SetSelected("Container Images")
	} else {
		radio.SetSelected("Local Tools")
	}

	middle := container.New(layout.NewVBoxLayout(), heading, label1, radio)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
