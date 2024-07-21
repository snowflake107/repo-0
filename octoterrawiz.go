package main

import (
	"github.com/mcasperson/OctoterraWizard/internal/steps"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

func main() {
	wiz := wizard.NewWizard()
	wiz.ShowWizardStep(steps.WelcomeStep{})
	wiz.Window.ShowAndRun()
}
