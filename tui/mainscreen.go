package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
)

type GroupMapEntry struct {
  Index   int
}

type Mainscreen struct {
  T      *TUI
  Canvas *tview.Grid

  Groups *tview.List
  Articles *tview.List

  CurrentGroupSelected int
  CurrentArticleSelected int

  GroupsMap map[string]GroupMapEntry
  GroupsList []string

  ArticlesDatasource   *[]models.Article
  ArticlesList []*models.Article
}

func(t *TUI) NewMainscreen(articlesDatasource *[]models.Article) (*Mainscreen) {
  mainscreen := new(Mainscreen)
  mainscreen.T = t

  mainscreen.ArticlesDatasource = articlesDatasource

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

  mainscreen.Groups = tview.NewList().
    SetWrapAround(false).
    ShowSecondaryText(false).
    SetHighlightFullLine(true).
    SetSelectedBackgroundColor(tcell.ColorTeal).
    SetSecondaryTextColor(tcell.ColorGrey).
    SetChangedFunc(mainscreen.changeHandler("group")).
    SetSelectedFunc(mainscreen.selectHandler("group"))

  mainscreen.Articles = tview.NewList().
    SetWrapAround(true).
    ShowSecondaryText(true).
    SetHighlightFullLine(true).
    SetSelectedBackgroundColor(tcell.ColorTeal).
    SetSecondaryTextColor(tcell.ColorGrey).
    SetChangedFunc(mainscreen.changeHandler("article")).
    SetSelectedFunc(mainscreen.selectHandler("article"))

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

func (mainscreen *Mainscreen) GetDefaultFocus() (tview.Primitive) {
  return mainscreen.Articles
}

func(mainscreen *Mainscreen) Refresh() {
  selectedGroup := mainscreen.CurrentGroupSelected
  selectedArticle := mainscreen.CurrentArticleSelected

  previousGroupsList := mainscreen.GroupsList
  mainscreen.GroupsList = []string{}
  // previousGroupsMap := mainscreen.GroupsMap
  mainscreen.GroupsMap = make(map[string]GroupMapEntry)
  mainscreen.Groups.Clear()

  mainscreen.ArticlesList = []*models.Article{}
  mainscreen.Articles.Clear()

  mainscreen.GroupsList = append(mainscreen.GroupsList, "*")
  mainscreen.GroupsMap["*"] = GroupMapEntry{
    Index: 0,
  }

  for _, article := range *mainscreen.ArticlesDatasource {
    if selectedGroup == 0 ||
      (selectedGroup != 0 &&
        article.Newsgroup == previousGroupsList[selectedGroup]) {
      mainscreen.Articles.AddItem(article.Subject, article.From, 0, nil)
      mainscreen.ArticlesList = append(mainscreen.ArticlesList, &article)
    }

    if _, ok := mainscreen.GroupsMap[article.Newsgroup]; !ok {
      mainscreen.GroupsList = append(mainscreen.GroupsList, article.Newsgroup)
      mainscreen.GroupsMap[article.Newsgroup] = GroupMapEntry{
        Index: (mainscreen.Groups.GetItemCount() - 1),
      }
    }
  }

  sort.Strings(mainscreen.GroupsList)
  for idx, group := range mainscreen.GroupsList {
    mainscreen.GroupsMap[group] = GroupMapEntry{
      Index: idx,
    }
    mainscreen.Groups.AddItem(group, "", 0, nil)
  }

  mainscreen.Groups.SetCurrentItem(selectedGroup)
  mainscreen.Articles.SetCurrentItem(selectedArticle)
  mainscreen.T.App.SetFocus(mainscreen.Articles)
}

func (mainscreen *Mainscreen) HandleInput(event *tcell.EventKey) (*tcell.EventKey) {
  switch event.Key() {
  case tcell.KeyCtrlG:
    mainscreen.T.App.SetFocus(mainscreen.Groups)
    return nil
  case tcell.KeyCtrlL:
    mainscreen.T.App.SetFocus(mainscreen.Articles)
    return nil
  case tcell.KeyRune:
    switch event.Rune() {
    case 'r':
      mainscreen.replyToArticle(mainscreen.ArticlesList[mainscreen.CurrentArticleSelected])
    }
    return nil
  }

  return event
}

func(mainscreen *Mainscreen) changeHandler(item string)(func(int, string, string, rune)) {
  return func(index int, text string, secondaryText string, shortcut rune) {
    switch(item) {
    case "group":
      mainscreen.CurrentGroupSelected = index
    case "article":
      mainscreen.CurrentArticleSelected = index
    }
  }
}

func(mainscreen *Mainscreen) selectHandler(item string)(func(int, string, string, rune)) {
  return func(index int, text string, secondaryText string, shortcut rune) {
    switch(item) {
    case "group":
      mainscreen.Refresh()
    case "article":
      OpenArticle(mainscreen.T.App, mainscreen.ArticlesList[index])
    }
  }
}

func(mainscreen *Mainscreen) replyToArticle(article *models.Article) {
  newArticle := models.NewArticle()

  newArticle.Subject = fmt.Sprintf("Re: %s", article.Subject)
  newArticle.InReplyToID = article.ID
  newArticle.Newsgroup = article.Newsgroup
  // TODO: newArticle.From =
  // TODO: newArticle.Organisation =
  newArticle.Body = fmt.Sprintf("\nOn %s %s wrote:\n> %s", time.Unix(0, article.Date * int64(time.Millisecond)).Format("Mon Jan _2 15:04:05 2006"), article.From, strings.Replace(article.Body, "\n", "\n> ", -1))

  _, err := OpenArticle(mainscreen.T.App, newArticle)
  if err != nil {
    // TODO: Show modal
    return
  }

  // TODO: Write reply
  mainscreen.T.ShowModal(
    "Do you want to post this article?",
    map[string]ModalButton{
      "[Y]es": {
        Rune: 'y',
        Callback: func() {
          return
        },
      },
      "[N]o": {
        Rune: 'n',
        Callback: func() {
          return
        },
      },
    })
}

