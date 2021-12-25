package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
)


type Mainscreen struct {
  Canvas *tview.Grid

  Groups *tview.List
  Articles *tview.List

  ArticlesDatasource   *[]models.Article
}

func(t *TUI) NewMainscreen(articlesDatasource *[]models.Article) (*Mainscreen) {
  mainscreen := new(Mainscreen)

  mainscreen.ArticlesDatasource = articlesDatasource

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

  mainscreen.Groups = tview.NewList().
    SetWrapAround(false).
    ShowSecondaryText(false).
    AddItem("List item 1", "Some explanatory text", 0, nil).
    AddItem("List item 2", "Some explanatory text", 0, nil).
    AddItem("List item 3", "Some explanatory text", 0, nil).
    AddItem("List item 4", "Some explanatory text", 0, nil).
    AddItem("Quit", "Press to exit", 'q', func() {
    })

  mainscreen.Articles = tview.NewList().
    SetWrapAround(true).
    SetHighlightFullLine(true).
    SetSelectedBackgroundColor(tcell.ColorTeal).
    SetSecondaryTextColor(tcell.ColorGrey).

    AddItem("List item 1", "Some explanatory text", 0, nil).
    AddItem("List item 2", "Some explanatory text", 0, nil).
    AddItem("List item 3", "Some explanatory text", 0, nil).
    AddItem("List item 4", "Some explanatory text", 0, nil)

	mainscreen.Canvas = tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 2, 0, 0, false)

	mainscreen.Canvas.AddItem(mainscreen.Groups, 1, 0, 1, 1, 0, 0, false).
		AddItem(mainscreen.Articles, 1, 1, 1, 1, 0, 0, false)

  return mainscreen
}

func (mainscreen *Mainscreen) GetCanvas() (tview.Primitive) {
  return mainscreen.Canvas
}

func(mainscreen *Mainscreen) Refresh() {
  addedGroups := make(map[string]bool)
  mainscreen.Groups.Clear()
  mainscreen.Articles.Clear()

  mainscreen.Groups.AddItem("*", "", 0, nil)

  for _, article := range *mainscreen.ArticlesDatasource {
    mainscreen.Articles.AddItem(article.Subject, article.From, 0, nil)

    if addedGroups[article.Newsgroup] != true {
      mainscreen.Groups.AddItem(article.Newsgroup, "", 0, nil)
      addedGroups[article.Newsgroup] = true
    }
  }

}

