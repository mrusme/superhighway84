package tui

import (
	"embed"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
  App                        *tview.Application
  Views                      map[string]View
  ActiveView                 string
}

type View interface {
  GetCanvas() (tview.Primitive)
  Draw()
}

func Init(embedfs *embed.FS) (*TUI) {
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

  logoBytes, err := embedfs.ReadFile("superhighway84.png")
  if err != nil {
    log.Panicln(err)
  }

  t.Views = make(map[string]View)
  t.Views["splashscreen"] = t.NewSplashscreen(&logoBytes)

  t.initInput()
  return t
}

func (t *TUI) initInput() {
	t.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
      t.Draw()
			return nil
		case tcell.KeyCtrlQ:
			t.App.Stop()
		}

		return event
	})
}

func(t *TUI) SetView(name string) {
  t.App.SetRoot(t.Views[name].GetCanvas(), true)
  t.ActiveView = name
}

func (t *TUI) Draw() {
  t.Views[t.ActiveView].Draw()
}

