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
[darkcyan] /  / / // [-][hotpink]/___/\_,_/ .__/\__/_/ /_//_/_/\_, /_//_/__,__/\_,_/\_, /\___/ /_/[-] [dimgray]%s[-]
[hotpink]                   /_/                  /___/                /___/[-]           [yellow]%s[-]
`

var STATS_TEMPLATE =
`[gray]⦿ %d PEERS[-]
[yellow]▲ %.2f[-] [gray]MB/s[-]
[teal]▼ %.2f[-] [gray]MB/s[-]
[yellow]▲ %.2f[-] [gray]MB[-]
[teal]▼ %.2f[-] [gray]MB[-]
`

var INFO_TEMPLATE = "%s"
const (
  COLOR_SUBJECT_UNREAD = "teal"
  COLOR_SUBJECT_READ   = "white"
)

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
  Preview  *tview.TextView

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
    SetWrapAround(false).
    ShowSecondaryText(false).
    SetHighlightFullLine(true).
    SetMainTextColor(tcell.ColorWhite).
    SetSelectedBackgroundColor(tcell.ColorHotPink).
    SetSelectedTextColor(tcell.ColorWhite).
    SetSecondaryTextColor(tcell.ColorGray).
    SetChangedFunc(mainscreen.changeHandler("group")).
    SetSelectedFunc(mainscreen.selectHandler("group"))
  mainscreen.Groups.
    SetBorder(true).
    SetBorderAttributes(tcell.AttrNone).
    SetBorderColor(tcell.ColorTeal)

  mainscreen.Articles = tview.NewList().
    SetWrapAround(false).
    ShowSecondaryText(true).
    SetHighlightFullLine(true).
    SetMainTextColor(tcell.ColorTeal).
    SetSelectedBackgroundColor(tcell.ColorHotPink).
    SetSelectedTextColor(tcell.ColorWhite).
    SetSecondaryTextColor(tcell.ColorGray).
    SetChangedFunc(mainscreen.changeHandler("article")).
    SetSelectedFunc(mainscreen.selectHandler("article"))
  mainscreen.Articles.
    SetBorder(true).
    SetBorderAttributes(tcell.AttrNone).
    SetBorderColor(tcell.ColorTeal)

  mainscreen.Preview = tview.NewTextView().
    SetText("").
    SetTextColor(tcell.ColorWhite).
    SetDynamicColors(true)
  mainscreen.Preview.
    SetBorder(true).
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
		SetRows(5, 0, 0, 1).
		SetColumns(30, 0, 14).
		SetBorders(false).
		AddItem(mainscreen.Header, 0, 0, 1, 2, 0, 0, false).
    AddItem(mainscreen.Stats,  0, 2, 1, 1, 0, 0, false).
		AddItem(mainscreen.Info,   3, 0, 1, 1, 0, 0, false).
		AddItem(mainscreen.Footer, 3, 1, 1, 2, 0, 0, false)

	mainscreen.Canvas.
    AddItem(mainscreen.Groups,   1, 0, 2, 1, 0, 0, false).
		AddItem(mainscreen.Articles, 1, 1, 1, 2, 0, 0, false).
		AddItem(mainscreen.Preview,  2, 1, 1, 2, 0, 0, false)

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

func (mainscreen *Mainscreen) SetVersion(version string, versionLatest string) {
  v := version
  if version == "v0.0.0" {
    v = "DeLorean @ 1.21 Gigawatts"
  }

  l := ""
  if versionLatest != version &&
    version != "v0.0.0" {
    l = fmt.Sprintf("%s update available!", versionLatest)
  }

  mainscreen.Header.SetText(
    fmt.Sprintf(HEADER_LOGO,
      v,
      l,
    ),
  )
}

func (mainscreen *Mainscreen) GetCanvas() (tview.Primitive) {
  return mainscreen.Canvas
}

func (mainscreen *Mainscreen) GetDefaultFocus() (tview.Primitive) {
  return mainscreen.Articles
}

func(mainscreen *Mainscreen) addNodeToArticlesList(level int, articlesNode *[]*models.Article, selectedGroup int, previousGroupsList []string) {
  // fmt.Fprintf(os.Stderr, "%s Node has %d items\n", strings.Repeat(" ", level * 3), len(*articlesNode))

  for i := 0; i < len(*articlesNode); i++ {
    article := (*articlesNode)[i]

    // fmt.Fprintf(os.Stderr, "%s   Item has ID %s and is in reply of ID %s and has %d replies\n", strings.Repeat(" ", level * 3), article.ID, article.InReplyToID, len(article.Replies))

    if selectedGroup == 0 ||
      (selectedGroup != 0 &&
        article.Newsgroup == previousGroupsList[selectedGroup]) {

      prefix := ""
      if level > 0 {
        if i < (len(*articlesNode) - 1) || len(article.Replies) > 0 {
          prefix = "[gray]├[-]"
        } else {
          prefix = "[gray]└[-]"
        }
      }

      prefixSub := " "
      if len(article.Replies) > 0 || (level > 0 && i < (len(*articlesNode) - 1)) {
        prefixSub = "[gray]│[-]"
      }

      subjectColor := COLOR_SUBJECT_UNREAD
      if article.Read == true {
        subjectColor = COLOR_SUBJECT_READ
      }

      mainscreen.Articles.AddItem(
        fmt.Sprintf(
          "%s%s[%s]%s[-]",
          prefix,
          strings.Repeat(" ", level),
          subjectColor,
          article.Subject,
        ),
        fmt.Sprintf(
          "%s%s in %s by [darkgray]%s[-]",
          prefixSub,
          strings.Repeat(" ", level),
          article.Newsgroup,
          article.From,
        ), 0, nil)
      mainscreen.ArticlesList = append(mainscreen.ArticlesList, article)

      if len(article.Replies) > 0 {
        mainscreen.addNodeToArticlesList((level + 1), &article.Replies, selectedGroup, previousGroupsList)
      }
    }

    if _, ok := mainscreen.GroupsMap[article.Newsgroup]; !ok {
      mainscreen.GroupsList = append(mainscreen.GroupsList, article.Newsgroup)
      mainscreen.GroupsMap[article.Newsgroup] = GroupMapEntry{
        Index: 0,
      }
    }
  }
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

  mainscreen.addNodeToArticlesList(0, mainscreen.T.ArticlesRoots, selectedGroup, previousGroupsList)

  sort.Strings(mainscreen.GroupsList)
  for idx, group := range mainscreen.GroupsList {
    mainscreen.GroupsMap[group] = GroupMapEntry{
      Index: idx,
    }
    mainscreen.Groups.AddItem(group, "", 0, nil)
  }

  mainscreen.Groups.SetCurrentItem(selectedGroup)
  mainscreen.Articles.SetCurrentItem(selectedArticle)

  mainscreen.changeHandler("group")(selectedGroup, "", "", 0)
  mainscreen.changeHandler("article")(selectedArticle, "", "", 0)

  mainscreen.T.App.SetFocus(mainscreen.Articles)
}

func (mainscreen *Mainscreen) HandleInput(event *tcell.EventKey) (*tcell.EventKey) {
  switch event.Key() {
  case tcell.KeyCtrlH:
    mainscreen.T.App.SetFocus(mainscreen.Groups)
    return nil
  case tcell.KeyCtrlJ:
    mainscreen.T.App.SetFocus(mainscreen.Preview)
    return nil
  case tcell.KeyCtrlL, tcell.KeyCtrlK:
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
    case 'g':
      if event.Rune() == 'G' {
        mainscreen.T.App.QueueEvent(tcell.NewEventKey(tcell.KeyEnd, 0, tcell.ModNone))
      } else {
        mainscreen.T.App.QueueEvent(tcell.NewEventKey(tcell.KeyHome, 0, tcell.ModNone))
      }
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
      if index < 0 || index >= len(mainscreen.ArticlesList) {
        return
      }
      mainscreen.CurrentArticleSelected = index
      mainscreen.renderPreview(mainscreen.ArticlesList[index])
      mainscreen.markAsRead(index, mainscreen.ArticlesList[index])
    }
  }
}

func(mainscreen *Mainscreen) selectHandler(item string)(func(int, string, string, rune)) {
  return func(index int, text string, secondaryText string, shortcut rune) {
    switch(item) {
    case "group":
      mainscreen.Refresh()
    case "article":
      mainscreen.T.OpenArticle(mainscreen.ArticlesList[index], true)
    }
  }
}

func(mainscreen *Mainscreen) renderPreview(article *models.Article) {
  mainscreen.Preview.SetText(fmt.Sprintf(
    "[gray]Date:[-] [darkgray]%s[-]\n[gray]Newsgroup:[-] [darkgray]%s[-]\n\n\n%s",
    MillisecondsToDate(article.Date),
    article.Newsgroup,
    article.Body,
  ))
  mainscreen.Preview.ScrollToBeginning()
}

func(mainscreen *Mainscreen) markAsRead(index int, article *models.Article) {
  article.Read = true
  mainText, secondaryText := mainscreen.Articles.GetItemText(index)
  updatedMainText := strings.Replace(
    mainText,
    fmt.Sprintf("[%s]", COLOR_SUBJECT_UNREAD),
    fmt.Sprintf("[%s]", COLOR_SUBJECT_READ),
    1,
  )
  mainscreen.Articles.SetItemText(index, updatedMainText, secondaryText)
  mainscreen.T.Cache.StoreArticle(article)
}

func(mainscreen *Mainscreen) submitNewArticle(group string) {
  newArticle := models.NewArticle()

  newArticle.Subject = ""
  newArticle.Newsgroup = group
  newArticle.From = mainscreen.T.Config.Profile.From
  newArticle.Organization = mainscreen.T.Config.Profile.Organization
  newArticle.Body = ""

  updatedNewArticle, err := mainscreen.T.OpenArticle(newArticle, false)
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
  newArticle.From = mainscreen.T.Config.Profile.From
  newArticle.Organization = mainscreen.T.Config.Profile.Organization
  newArticle.Body = fmt.Sprintf("\nOn %s %s wrote:\n> %s", MillisecondsToDate(article.Date), article.From, strings.Replace(article.Body, "\n", "\n> ", -1))

  updatedNewArticle, err := mainscreen.T.OpenArticle(newArticle, false)
  if err != nil {
    mainscreen.T.ShowErrorModal(err.Error())
    return
  }

  if strings.TrimSpace(updatedNewArticle.Body) == strings.TrimSpace(newArticle.Body) {
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

