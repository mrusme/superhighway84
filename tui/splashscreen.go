package tui

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/rivo/tview"
)


type Splashscreen struct {
  Canvas *tview.TextView
  ImageBytes []byte
}

func(t *TUI) NewSplashscreen(logo *[]byte) (*Splashscreen) {
  splashscreen := new(Splashscreen)

  canvas := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true)
  canvas.SetBorder(false)
  canvas.Clear()

  splashscreen.ImageBytes = *logo

  splashscreen.Canvas = canvas
  return splashscreen
}

func (splashscreen *Splashscreen) GetCanvas() (tview.Primitive) {
  return splashscreen.Canvas
}

func(splashscreen *Splashscreen) Refresh() {
  _, _, w, h := splashscreen.Canvas.Box.GetRect()

  // TODO:
  // (h * 2) is a workaround for what looks like a bug in
  // https://github.com/eliukblau/pixterm/blob/master/pkg/ansimage/ansimage.go
  // Depending on the dithering setting the h/w changes significantly.
  logoImage, err := ansimage.NewScaledFromReader(bytes.NewReader(splashscreen.ImageBytes), (h * 2), w, color.Black, ansimage.ScaleModeFill, ansimage.NoDithering)
  if err != nil {
    return
  }
  // splashscreen.Canvas.Clear()
	fmt.Fprint(splashscreen.Canvas, tview.TranslateANSI(logoImage.RenderExt(false, false)))
}

