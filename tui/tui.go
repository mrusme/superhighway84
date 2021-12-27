package tui

import (
	"embed"
	"log"
	"time"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/mrusme/superhighway84/config"
	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
	"go.uber.org/zap"
)

type TUI struct {
  App                        *tview.Application
  Views                      map[string]View
  ActiveView                 string

  Modal                      *tview.Modal
  ModalVisible               bool
  ModalButtons               map[string]ModalButton

  ArticlesDatasource         *[]models.Article

  CallbackRefreshArticles    func() (error)
  CallbackSubmitArticle      func(article *models.Article) (error)

  Config                     *config.Config
  Logger                     *zap.Logger

  Stats                      map[string]int64
}

type View interface {
  GetCanvas() (tview.Primitive)
  GetDefaultFocus() (tview.Primitive)

  Refresh()

  HandleInput(event *tcell.EventKey) (*tcell.EventKey)
}

type ModalButton struct {
  Rune      rune
  Callback  func()
}

func Init(embedfs *embed.FS, cfg *config.Config, logger *zap.Logger) (*TUI) {
  t := new(TUI)

  tview.Styles = tview.Theme{
    PrimitiveBackgroundColor:    tcell.ColorDefault,
    ContrastBackgroundColor:     tcell.ColorPink,
    MoreContrastBackgroundColor: tcell.ColorTeal,
    BorderColor:                 tcell.ColorWhite,
    TitleColor:                  tcell.ColorWhite,
    GraphicsColor:               tcell.ColorWhite,
    PrimaryTextColor:            tcell.ColorDefault,
    SecondaryTextColor:          tcell.ColorBlue,
    TertiaryTextColor:           tcell.ColorGreen,
    InverseTextColor:            tcell.ColorBlack,
    ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
  }

  t.App = tview.NewApplication()
  t.Config = cfg
  t.Logger = logger

  t.Stats = make(map[string]int64)

  logoBytes, err := embedfs.ReadFile("superhighway84.jpeg")
  if err != nil {
    log.Panicln(err)
  }

  t.Views = make(map[string]View)
  t.Views["splashscreen"] = t.NewSplashscreen(&logoBytes)
  t.Views["mainscreen"] = t.NewMainscreen()

  t.ModalVisible = false

  t.initInput()
  return t
}

func (t *TUI) initInput() {
	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
      t.RefreshMainscreen()
      t.SetInfo(true)
			return nil
		case tcell.KeyCtrlQ:
			t.App.Stop()
      return nil
    default:
      if t.ModalVisible == true {
        for _, modalButton := range t.ModalButtons {
          if modalButton.Rune == '*' ||
            (event.Key() == tcell.KeyRune &&
              unicode.ToLower(modalButton.Rune) == unicode.ToLower(event.Rune())) {
            modalButton.Callback()
            t.HideModal()
            return nil
          }
        }
        return nil
      } else {
        return t.Views[t.ActiveView].HandleInput(event)
      }
		}
	})
}

func (t *TUI) Launch() {
  go func() {
    time.Sleep(time.Millisecond * 200)
    t.SetView("splashscreen", true)
    t.Refresh()
    t.App.Draw()
  }()

  go func() {
    for {
      t.RefreshMainscreen()
      t.SetInfo(true)
      time.Sleep(time.Second * 60)
    }
  }()

  if err := t.App.Run(); err != nil {
    panic(err)
  }
}

func(t *TUI) SetView(name string, redraw bool) {
  t.ActiveView = name

  t.App.SetRoot(t.Views[t.ActiveView].GetCanvas(), true).
    SetFocus(t.Views[t.ActiveView].GetDefaultFocus())

  if redraw {
    t.App.Draw()
  }
}

func (t *TUI) RefreshMainscreen() {
  if t.ActiveView != "mainscreen" {
    return
  }
  t.RefreshData()
  t.Refresh()
  return
}

func (t *TUI) Refresh() {
  t.Views[t.ActiveView].Refresh()
  return
}

func (t *TUI) RefreshData() {
  if t.CallbackRefreshArticles != nil {
    err := t.CallbackRefreshArticles()
    if err != nil {
      t.ShowErrorModal(err.Error())
    }
  }
}

func(t *TUI) ShowModal(text string, buttons map[string]ModalButton) {
  t.Modal = tview.NewModal().
    SetTextColor(tcell.ColorBlack).
    SetBackgroundColor(tcell.ColorTeal).
    SetButtonBackgroundColor(tcell.ColorHotPink).
    SetButtonTextColor(tcell.ColorWhite).
    SetText(text)
    // SetDoneFunc(func(buttonIndex int, buttonLabel string) {
    //   modalButton := buttons[buttonLabel]
    //   modalButton.Callback()
    //   t.HideModal()
    // })
  var buttonLabels []string
  for buttonLabel := range buttons {
    buttonLabels = append(buttonLabels, buttonLabel)
  }
  t.Modal.AddButtons(buttonLabels)

  t.ModalVisible = true
  t.ModalButtons = buttons
  t.App.SetRoot(t.Modal, false).SetFocus(t.Modal)
}

func(t *TUI) HideModal() {
  t.ModalVisible = false
  t.SetView(t.ActiveView, false)
}

func(t *TUI) ShowErrorModal(text string) {
  t.ShowModal(
    text,
    map[string]ModalButton{
      "(F)uck": {
        Rune: '*',
        Callback: func() {
          return
        },
      },
    })
}

func (t *TUI) SetInfo(refresh bool) {
  if refresh == true {
    t.Views["mainscreen"].(*Mainscreen).SetInfo(map[string]string{
      "refresh": "[green]◙[-]",
    })
    time.AfterFunc(time.Second * 3, func() {
      t.SetInfo(false)
    })
  } else {
    t.Views["mainscreen"].(*Mainscreen).SetInfo(map[string]string{
      "refresh": "[grey]◙[-]",
    })
  }
}

func (t* TUI) SetStats(peers, rateIn, rateOut, totalIn, totalOut int64) () {
  t.Stats["peers"] = peers
  t.Stats["rate_in"] = rateIn
  t.Stats["rate_out"] = rateOut
  t.Stats["total_in"] = totalIn
  t.Stats["total_out"] = totalOut

  t.Views["mainscreen"].(*Mainscreen).SetStats(t.Stats)
  t.App.Draw()
}

