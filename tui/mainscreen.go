package tui

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
)

var HEADER_LOGO =
`[white]    _  _ _ __ [-][hotpink]____                  __   _      __                   ___  ____[-]
[teal]   /  / / // [-][hotpink]/ __/_ _____  ___ ____/ /  (_)__ _/ / _    _____ ___ __( _ )/ / /[-]
[teal]  _\ _\_\_\\_[-][fuchsia]\ \/ // / _ \/ -_) __/ _ \/ / _ \/ _ \ |/|/ / _ \/ // / _  /_  _/[-]
[darkcyan] /  / / // [-][hotpink]/___/\_,_/ .__/\__/_/ /_//_/_/\_, /_//_/__,__/\_,_/\_, /\___/ /_/[-]
[hotpink]                   /_/                  /___/                /___/[-]
`

var STATS_TEMPLATE =
`[grey]⦿ %d PEERS[-]
[yellow]▲ %.2f[-] [grey]MB/s[-]
[teal]▼ %.2f[-] [grey]MB/s[-]
[yellow]▲ %.2f[-] [grey]MB[-]
[teal]▼ %.2f[-] [grey]MB[-]
`

var INFO_TEMPLATE = "%s"

type GroupMapEntry struct {
  Index   int
}

type Mainscreen struct {
  T      *TUI
  Canvas *tview.Grid

  Header *tview.TextView
  Stats  *tview.TextView
  Info   *tview.TextView
  Footer *tview.TextView

  Groups *tview.List
  Articles *tview.List

  CurrentGroupSelected int
  CurrentArticleSelected int

  GroupsMap map[string]GroupMapEntry
  GroupsList []string

  ArticlesList []*models.Article
}

func(t *TUI) NewMainscreen() (*Mainscreen) {
  mainscreen := new(Mainscreen)
  mainscreen.T = t

  mainscreen.Groups = tview.NewList().
    SetWrapAround(true).
    ShowSecondaryText(false).
    SetHighlightFullLine(true).
    SetSelectedBackgroundColor(tcell.ColorHotPink).
    SetSelectedTextColor(tcell.ColorWhite).
    SetSecondaryTextColor(tcell.ColorGrey).
    SetChangedFunc(mainscreen.changeHandler("group")).
    SetSelectedFunc(mainscreen.selectHandler("group"))
  mainscreen.Groups.
    SetBorder(true).
    SetBorderAttributes(tcell.AttrNone).
    SetBorderColor(tcell.ColorTeal)

  mainscreen.Articles = tview.NewList().
    SetWrapAround(true).
    ShowSecondaryText(true).
    SetHighlightFullLine(true).
    SetSelectedBackgroundColor(tcell.ColorHotPink).
    SetSelectedTextColor(tcell.ColorWhite).
    SetSecondaryTextColor(tcell.ColorGrey).
    SetChangedFunc(mainscreen.changeHandler("article")).
    SetSelectedFunc(mainscreen.selectHandler("article"))
  mainscreen.Articles.
    SetBorder(true).
    SetBorderAttributes(tcell.AttrNone).
    SetBorderColor(tcell.ColorTeal)

  mainscreen.Header = tview.NewTextView().
    SetText(HEADER_LOGO).
    SetTextColor(tcell.ColorHotPink).
    SetDynamicColors(true)
  mainscreen.Header.SetBorder(false)

  mainscreen.Stats = tview.NewTextView().
    SetText("").
    SetTextColor(tcell.ColorHotPink).
    SetDynamicColors(true)
  mainscreen.Stats.SetBorder(false)

  mainscreen.Info = tview.NewTextView().
    SetText("").
    SetTextColor(tcell.ColorHotPink).
    SetDynamicColors(true)
  mainscreen.Info.SetBorder(false).
    SetBorderPadding(0, 0, 1, 1)

  mainscreen.Footer = tview.NewTextView().
    SetText("It really whips the llama's ass").
    SetTextColor(tcell.ColorHotPink).
    SetTextAlign(tview.AlignRight)
  mainscreen.Footer.SetBorder(false).
    SetBorderPadding(0, 0, 1, 1)

	mainscreen.Canvas = tview.NewGrid().
		SetRows(5, 0, 1).
		SetColumns(30, 0, 14).
		SetBorders(false).
		AddItem(mainscreen.Header, 0, 0, 1, 2, 0, 0, false).
    AddItem(mainscreen.Stats,  0, 2, 1, 1, 0, 0, false).
		AddItem(mainscreen.Info,   2, 0, 1, 1, 0, 0, false).
		AddItem(mainscreen.Footer, 2, 1, 1, 2, 0, 0, false)

	mainscreen.Canvas.
    AddItem(mainscreen.Groups,   1, 0, 1, 1, 0, 0, false).
		AddItem(mainscreen.Articles, 1, 1, 1, 2, 0, 0, false)

  return mainscreen
}

func (mainscreen *Mainscreen) SetFooter(text string) {
  mainscreen.Footer.SetText(text)
}

func (mainscreen *Mainscreen) SetStats(stats map[string]int64) {
  peers    := stats["peers"]
  totalIn  := float64(stats["total_in"])  / 1024.0 / 1024.0
  totalOut := float64(stats["total_out"]) / 1024.0 / 1024.0
  rateIn   := float64(stats["rate_in"])   / 1024.0 / 1024.0
  rateOut  := float64(stats["rate_out"])  / 1024.0 / 1024.0

  mainscreen.Stats.SetText(
    fmt.Sprintf(STATS_TEMPLATE,
      peers,
      rateOut,
      rateIn,
      totalOut,
      totalIn,
    ),
  )
}

func (mainscreen *Mainscreen) SetInfo(info map[string]string) {
  refresh := info["refresh"]
  mainscreen.Info.SetText(
    fmt.Sprintf(INFO_TEMPLATE,
      refresh,
    ),
  )
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

  for i := 0; i < len(*mainscreen.T.ArticlesDatasource); i++ {
    article := (*mainscreen.T.ArticlesDatasource)[i]
    if selectedGroup == 0 ||
      (selectedGroup != 0 &&
        article.Newsgroup == previousGroupsList[selectedGroup]) {
      mainscreen.Articles.AddItem(fmt.Sprintf("[teal]%s[-]", article.Subject), fmt.Sprintf("On [lightgray]%s[-] by %s", MillisecondsToDate(article.Date), article.From), 0, nil)
      mainscreen.ArticlesList = append(mainscreen.ArticlesList, &article)
    }

    if _, ok := mainscreen.GroupsMap[article.Newsgroup]; !ok {
      mainscreen.GroupsList = append(mainscreen.GroupsList, article.Newsgroup)
      mainscreen.GroupsMap[article.Newsgroup] = GroupMapEntry{
        Index: 0,
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
  case tcell.KeyCtrlH:
    mainscreen.T.App.SetFocus(mainscreen.Groups)
    return nil
  case tcell.KeyCtrlL:
    mainscreen.T.App.SetFocus(mainscreen.Articles)
    return nil
  case tcell.KeyRune:
    switch unicode.ToLower(event.Rune()) {
    case 'n':
      mainscreen.submitNewArticle(mainscreen.GroupsList[mainscreen.CurrentGroupSelected])
      return nil
    case 'r':
      mainscreen.replyToArticle(mainscreen.ArticlesList[mainscreen.CurrentArticleSelected])
      return nil
    case 'j':
       mainscreen.T.App.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
       return nil
    case 'k':
       mainscreen.T.App.QueueEvent(tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone))
       return nil
    case 'h':
       mainscreen.T.App.QueueEvent(tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone))
       return nil
    case 'l':
       mainscreen.T.App.QueueEvent(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone))
       return nil
    }
    return event
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
      mainscreen.T.OpenArticle(mainscreen.ArticlesList[index])
    }
  }
}

func(mainscreen *Mainscreen) submitNewArticle(group string) {
  newArticle := models.NewArticle()

  newArticle.Subject = ""
  newArticle.Newsgroup = ""
  newArticle.From = "nya"
  newArticle.Organization = "su"
  newArticle.Body = ""

  updatedNewArticle, err := mainscreen.T.OpenArticle2(newArticle)
  if err != nil {
    mainscreen.T.ShowErrorModal(err.Error())
    return
  }

  mainscreen.T.ShowModal(
    "Do you want to submit this new article?",
    map[string]ModalButton{
      "(Y)es": {
        Rune: 'y',
        Callback: func() {
          if mainscreen.T.CallbackSubmitArticle != nil {
            mainscreen.T.CallbackSubmitArticle(&updatedNewArticle)
          }
          return
        },
      },
      "(N)ope": {
        Rune: 'n',
        Callback: func() {
          return
        },
      },
    })

}

func(mainscreen *Mainscreen) replyToArticle(article *models.Article) {
  newArticle := models.NewArticle()

  newArticle.Subject = fmt.Sprintf("Re: %s", article.Subject)
  newArticle.InReplyToID = article.ID
  newArticle.Newsgroup = article.Newsgroup
  newArticle.From = "nya"
  newArticle.Organization = "su"
  newArticle.Body = fmt.Sprintf("\nOn %s %s wrote:\n> %s", MillisecondsToDate(article.Date), article.From, strings.Replace(article.Body, "\n", "\n> ", -1))

  updatedNewArticle, err := mainscreen.T.OpenArticle2(newArticle)
  if err != nil {
    mainscreen.T.ShowErrorModal(err.Error())
    return
  }

  mainscreen.T.ShowModal(
    "Do you want to submit this reply?",
    map[string]ModalButton{
      "(Y)es": {
        Rune: 'y',
        Callback: func() {
          if mainscreen.T.CallbackSubmitArticle != nil {
            mainscreen.T.CallbackSubmitArticle(&updatedNewArticle)
          }
          return
        },
      },
      "(N)ope": {
        Rune: 'n',
        Callback: func() {
          return
        },
      },
    })
}

