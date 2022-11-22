package dcshmd

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dimchansky/dcs-hmd/ka50rotorpitch"
	"github.com/dimchansky/dcs-hmd/ka50rotorrpm"
	"github.com/dimchansky/dcs-hmd/utils"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600

	fontBaseSize = 18
	dpi          = 72

	rowWidth  = 20
	rowHeight = 20
)

var (
	textColor   = color.NRGBA{G: 0xff, A: 0xff}
	shadowColor = color.NRGBA{A: 0xff}
)

func NewHUD() (*HUD, error) {
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenFilterEnabled(false)

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetInitFocused(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	const (
		indicatorHeight = 400
		xSpan           = rowWidth / 2
		ySpan           = rowHeight
	)

	rotorPitchIndicator := ka50rotorpitch.NewIndicator(&ka50rotorpitch.IndicatorConfig{
		Width:           rowWidth * 3,
		Height:          indicatorHeight,
		TickLength:      rowWidth,
		MinorTickLength: rowWidth / 2,
		LineWidth:       2,
		Color:           textColor,
		BorderColor:     shadowColor,
		Rect: image.Rectangle{
			Min: image.Pt(xSpan, ySpan),
			Max: image.Pt(rowWidth*2+xSpan-1, indicatorHeight-ySpan-1),
		},
	})
	rotorRPMIndicator := ka50rotorrpm.NewIndicator(&ka50rotorrpm.IndicatorConfig{
		Width:           rowWidth * 3,
		Height:          indicatorHeight,
		TickLength:      rowWidth,
		MinorTickLength: rowWidth * 3 / 4,
		LineWidth:       2,
		Color:           textColor,
		BorderColor:     shadowColor,
		Rect: image.Rectangle{
			Min: image.Pt(xSpan, ySpan),
			Max: image.Pt(rowWidth*2+xSpan-1, indicatorHeight-ySpan-1),
		},
	})

	ff, err := NewFontFace(fontBaseSize, dpi)
	if err != nil {
		return nil, err
	}

	hud := &HUD{
		fontFace:            ff,
		rotorPitchIndicator: rotorPitchIndicator,
		rotorRPMIndicator:   rotorRPMIndicator,
	}

	return hud, nil
}

type HUD struct {
	once sync.Once

	fontFace            *FontFace
	rotorPitchIndicator *ka50rotorpitch.Indicator
	rotorRPMIndicator   *ka50rotorrpm.Indicator

	rotorPitchImg redrawnImage
	rotorRPMImg   redrawnImage
}

func (h *HUD) Close() error {
	return fmt.Errorf("failed to close font face: %w", h.fontFace.Close())
}

func (h *HUD) Update() error {
	h.once.Do(enableCurrentProcessWindowClickThroughAsync)

	h.rotorPitchImg.Update(h.rotorPitchIndicator.GetImage())
	h.rotorRPMImg.Update(h.rotorRPMIndicator.GetImage())

	return nil
}

func (h *HUD) Draw(screen *ebiten.Image) {
	rotorPitchImg := &h.rotorPitchImg
	rotorRPMImg := &h.rotorRPMImg

	if !rotorPitchImg.NeedToDraw &&
		!rotorRPMImg.NeedToDraw {
		return
	}

	op := &ebiten.DrawImageOptions{CompositeMode: ebiten.CompositeModeCopy}
	op.GeoM.Translate(0, rowHeight)
	rotorPitchImg.DrawOn(screen, op)

	op.GeoM.Translate(float64(rotorPitchImg.Size().X), 0)
	rotorRPMImg.DrawOn(screen, op)
}

func (h *HUD) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// SetRotorPitch is thread-safe to update rotor pitch.
func (h *HUD) SetRotorPitch(val float64) {
	h.rotorPitchIndicator.SetRotorPitch(val)
}

// SetRotorRPM is thread-safe to update rotor RPM.
func (h *HUD) SetRotorRPM(val float64) {
	h.rotorRPMIndicator.SetRotorRPM(val)
}

func enableCurrentProcessWindowClickThroughAsync() {
	go utils.EnableCurrentProcessWindowClickThrough()
}

type redrawnImage struct {
	img        *ebiten.Image
	NeedToDraw bool
}

func (i *redrawnImage) Update(img *ebiten.Image, isRedrawn bool) {
	i.NeedToDraw = i.img == nil || isRedrawn
	i.img = img
}

func (i *redrawnImage) DrawOn(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if i.NeedToDraw {
		screen.DrawImage(i.img, op)
		i.NeedToDraw = false
	}
}
func (i *redrawnImage) Size() image.Point {
	return i.img.Bounds().Size()
}
