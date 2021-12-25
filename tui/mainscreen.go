package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


type Mainscreen struct {
  Canvas *tview.Grid
}

func(t *TUI) NewMainscreen() (*Mainscreen) {
  mainscreen := new(Mainscreen)

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

  groups := tview.NewList().
    AddItem("List item 1", "Some explanatory text", 'a', nil).
    AddItem("List item 2", "Some explanatory text", 'b', nil).
    AddItem("List item 3", "Some explanatory text", 'c', nil).
    AddItem("List item 4", "Some explanatory text", 'd', nil).
    AddItem("Quit", "Press to exit", 'q', func() {
    })

  articles := tview.NewList().
    SetWrapAround(true).
    SetHighlightFullLine(true).
    SetSelectedBackgroundColor(tcell.ColorTeal).
    SetSecondaryTextColor(tcell.ColorGrey).

    AddItem("List item 1", "Some explanatory text", ' ', nil).
    AddItem("List item 2", "Some explanatory text", ' ', nil).
    AddItem("List item 3", "Some explanatory text", ' ', nil).
    AddItem("List item 4", "Some explanatory text", ' ', nil)

	mainscreen.Canvas = tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 2, 0, 0, false)

	mainscreen.Canvas.AddItem(groups, 1, 0, 1, 1, 0, 0, false).
		AddItem(articles, 1, 1, 1, 1, 0, 0, false)

  return mainscreen
}

func (mainscreen *Mainscreen) GetCanvas() (tview.Primitive) {
  return mainscreen.Canvas
}

func(mainscreen *Mainscreen) Refresh() {
}

