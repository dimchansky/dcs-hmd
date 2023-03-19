package dcshmd

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dimchansky/dcs-hmd/aircraft/ka-50/devices/rotorpitch"
	"github.com/dimchansky/dcs-hmd/aircraft/ka-50/devices/rotorrpm"
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

	rotorPitchIndicator := rotorpitch.NewIndicator(&rotorpitch.IndicatorConfig{
		Width:           rowWidth * 3,
		Height:          indicatorHeight,
		TickLength:      rowWidth,
		MinorTickLength: rowWidth / 2,
		LineWidth:       2,
		Color:           textColor,
		BorderColor:     shadowColor,
		Rect: image.Rectangle{
			Min: image.Pt(xSpan, ySpan),
			Max: image.Pt(rowWidth*2+xSpan, indicatorHeight-ySpan),
		},
	})
	rotorRPMIndicator := rotorrpm.NewIndicator(&rotorrpm.IndicatorConfig{
		Width:           rowWidth * 3,
		Height:          indicatorHeight,
		TickLength:      rowWidth,
		MinorTickLength: rowWidth * 3 / 4,
		LineWidth:       2,
		Color:           textColor,
		BorderColor:     shadowColor,
		Rect: image.Rectangle{
			Min: image.Pt(xSpan, ySpan),
			Max: image.Pt(rowWidth*2+xSpan, indicatorHeight-ySpan),
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
	rotorPitchIndicator *rotorpitch.Indicator
	rotorRPMIndicator   *rotorrpm.Indicator

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

// redrawnImage represents an image that needs to be drawn again.
type redrawnImage struct {
	img        *ebiten.Image
	NeedToDraw bool
}

// Update updates the image and sets the NeedToDraw flag based on whether the image has been redrawn or not,
// or whether NeedToDraw was already true. If img is nil, NeedToDraw will always be true.
// The NeedToDraw flag is used to indicate whether the image needs to be drawn on the screen or not.
// If the flag is true, the image will be redrawn on the screen in the next frame.
// This flag is set to false by the DrawOn method after the image is drawn on the screen.
func (i *redrawnImage) Update(img *ebiten.Image, isRedrawn bool) {
	i.NeedToDraw = i.NeedToDraw || i.img == nil || isRedrawn
	i.img = img
}

// DrawOn draws the image on the given screen with the given options if the NeedToDraw flag is true.
// After drawing, NeedToDraw is set to false.
func (i *redrawnImage) DrawOn(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if i.NeedToDraw {
		screen.DrawImage(i.img, op)
		i.NeedToDraw = false
	}
}

// Size returns the size of the image.
func (i *redrawnImage) Size() image.Point {
	return i.img.Bounds().Size()
}
