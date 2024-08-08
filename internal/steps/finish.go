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

	heading := widget.NewLabel("Completed")
	heading.TextStyle = fyne.TextStyle{Bold: true}

	intro := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		The migration tool has now completed.
		Note that there are resources and settings that can not be migrated.
		Where possible, these resources can be manually recreated or updated in the destination space:
		* Config-as-Code settings
		* Certificates
		* Account, feed, Git, ServiceNow, Jira, Sumo, and Slunk credentials
		* Sensitive variables defined directly in a step, for example in a step template that has a sensitive parameter, or
          steps like the IIS or Tomcat steps that directly expose sensitive fields
		* Tenant sensitive variables
		* Users and teams
		* Subscriptions
		* Audit logs
		* Deployment and runbook run history
		* Built-in feed packages
		* Build information
		* Email settings
	`))
	middle := container.New(layout.NewVBoxLayout(), heading, intro)

	content := container.NewBorder(nil, nil, nil, nil, middle)

	return content
}
