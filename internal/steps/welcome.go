package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type WelcomeStep struct {
}

func (s WelcomeStep) GetContainer() *fyne.Container {
	label1 := widget.NewLabel("Label 1")
	return container.New(layout.NewFormLayout(), label1)
}
