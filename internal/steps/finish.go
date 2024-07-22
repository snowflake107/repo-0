package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type FinishStep struct {
	BaseStep
	Wizard wizard.Wizard
}

func (s FinishStep) GetContainer(parent fyne.Window) *fyne.Container {

	intro := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Your Octopus space now has a new project called "Octoterra Space Management" that contains two runbooks:
		__ 1. Serialize Space, which serializes space level resources to a Terraform module
		__ 2. Deploy Space, which deploys the Terraform module to a new space
	`))
	intro2 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Run the __ 1. Serialize Space runbook first, then run the __ 2. Deploy Space runbook.
	`))
	intro3 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Each project now contains two runbooks:
		__ 1. Serialize Project, which serializes the project to a Terraform module
		__ 2. Deploy Project, which deploys the Terraform module to a new space
	`))
	intro4 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		For each project, run the __ 1. Serialize Project first, then run the __ 2. Deploy Project runbook.
	`))
	middle := container.New(layout.NewVBoxLayout(), intro, intro2, intro3, intro4)

	content := container.NewBorder(nil, nil, nil, nil, middle)

	return content
}
