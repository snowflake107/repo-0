package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"os/exec"
)

type TestTerraformStep struct {
	BaseStep
	Wizard wizard.Wizard
}

func (s TestTerraformStep) GetContainer() *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(WelcomeStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(OctopusDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})
	next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		You must have Terraform installed to use this tool.
		Click the "Test" button to check if Terraform is installed.
	`))
	result := widget.NewLabel("")
	testTerraform := widget.NewButton("Test", func() {
		cmd := exec.Command("terraform", "-version")
		if err := cmd.Run(); err != nil {
			result.SetText("Terraform does not appear to be installed. You must install Terraform before proceeding.")
		} else {
			result.SetText("Terraform is installed. Click the Next button to proceed.")
			next.Enable()
		}
	})
	middle := container.New(layout.NewVBoxLayout(), label1, testTerraform, result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
