package tui

import (
	"embed"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mrusme/superhighway84/models"
	"github.com/rivo/tview"
)

type TUI struct {
  App                        *tview.Application
  Views                      map[string]View
  ActiveView                 string
}

type View interface {
  GetCanvas() (tview.Primitive)
  Refresh()
}

func Init(embedfs *embed.FS, articlesDatasource *[]models.Article) (*TUI) {
  t := new(TUI)

  tview.Styles = tview.Theme{
    PrimitiveBackgroundColor:    tcell.ColorDefault,
    ContrastBackgroundColor:     tcell.ColorTeal,
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

  logoBytes, err := embedfs.ReadFile("superhighway84.jpeg")
  if err != nil {
    log.Panicln(err)
  }

  t.Views = make(map[string]View)
  t.Views["splashscreen"] = t.NewSplashscreen(&logoBytes)
  t.Views["mainscreen"] = t.NewMainscreen(articlesDatasource)

  t.initInput()
  return t
}

func (t *TUI) initInput() {
	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
      t.Refresh()
			return nil
		case tcell.KeyCtrlQ:
			t.App.Stop()
		}

		return event
	})
}

func (t *TUI) Launch() {
  go func() {
    time.Sleep(time.Millisecond * 200)
    t.SetView("splashscreen")
    t.Refresh()
    t.App.Draw()
  }()

  if err := t.App.Run(); err != nil {
    panic(err)
  }
}

func(t *TUI) SetView(name string) {
  t.ActiveView = name
  t.App.SetRoot(t.Views[t.ActiveView].GetCanvas(), true)
  t.App.Draw()
}

func (t *TUI) Refresh() {
  t.Views[t.ActiveView].Refresh()
}

