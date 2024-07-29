package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type BackendSelectionStep struct {
	BaseStep
	Wizard wizard.Wizard
}

var AwsS3 = "AWS S3"
var AzureStorage = "Azure Storage"

func (s BackendSelectionStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, _ := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(ToolsSelectionStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		if s.State.BackendType == AzureStorage {
			s.Wizard.ShowWizardStep(AzureTerraformStateStep{
				Wizard:   s.Wizard,
				BaseStep: BaseStep{State: s.State}})
		} else {
			s.Wizard.ShowWizardStep(AwsTerraformStateStep{
				Wizard:   s.Wizard,
				BaseStep: BaseStep{State: s.State}})
		}
	})

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Terraform requires a backend to manage its state. Select a backend from the list below.
	`))

	radio := widget.NewRadioGroup([]string{AwsS3, AzureStorage}, func(value string) {
		s.State.BackendType = value
	})

	if s.State.BackendType == "" {
		radio.SetSelected(AzureStorage)
	} else {
		radio.SetSelected(s.State.BackendType)
	}

	middle := container.New(layout.NewVBoxLayout(), label1, radio)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
