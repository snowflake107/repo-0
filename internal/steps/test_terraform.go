package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"net/url"
	"os/exec"
)

type TestTerraformStep struct {
	BaseStep
	Wizard wizard.Wizard
	result *widget.Label
}

func (s TestTerraformStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, _ := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(WelcomeStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.result.SetText("🔵 Testing Terraform installation.")

		if s.Execute() {
			s.Wizard.ShowWizardStep(OctopusDetails{
				Wizard:   s.Wizard,
				BaseStep: BaseStep{State: s.State}})
		} else {
			s.notInstalled()
		}
	})

	heading := widget.NewLabel("Test Terraform")
	heading.TextStyle = fyne.TextStyle{Bold: true}

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		You must have Terraform installed to use this tool.
		Click the "Next" button to check if Terraform is installed.
	`))

	linkUrl, _ := url.Parse("https://developer.hashicorp.com/terraform/install")
	link := widget.NewHyperlink("Learn how to install Terraform.", linkUrl)

	s.result = widget.NewLabel("")

	middle := container.New(layout.NewVBoxLayout(), heading, label1, link, s.result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s TestTerraformStep) notInstalled() {
	s.result.SetText("🔴 Terraform does not appear to be installed. You must install Terraform before proceeding.")
}

func (s TestTerraformStep) Execute() bool {
	cmd := exec.Command("terraform", "-version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
