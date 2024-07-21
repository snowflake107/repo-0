package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/state"
)

type BaseStep struct {
	State state.State
}

func (s BaseStep) BuildNavigation(previousCallback func(), nextCallback func()) (*fyne.Container, *widget.Button, *widget.Button) {
	previous := widget.NewButton("< Previous", previousCallback)
	next := widget.NewButton("Next >", nextCallback)
	bottom := container.New(layout.NewGridLayout(2), previous, next)
	return bottom, previous, next
}
