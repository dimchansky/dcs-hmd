package dcshmd

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dimchansky/dcs-hmd/ka50rotorpitch"
	"github.com/dimchansky/dcs-hmd/ka50rotorrpm"

	_ "github.com/silbinarywolf/preferdiscretegpu"
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
	ebiten.SetTPS(20)

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
		ySpan           = rowHeight / 2
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
	fontFace            *FontFace
	rotorPitchIndicator *ka50rotorpitch.Indicator
	rotorRPMIndicator   *ka50rotorrpm.Indicator

	rotorPitchImg redrawnImage
	rotorRPMImg   redrawnImage
}

func (h *HUD) Close() error {
	return h.fontFace.Close()
}

var (
	dPitch = 0.05
	dRPM   = 0.1
)

func (h *HUD) Update() error {
	rotorPitch := h.rotorPitchIndicator.GetRotorPitch()
	if rotorPitch >= 15 {
		dPitch = -0.1
	} else if rotorPitch <= 1 {
		dPitch = 0.1
	}
	h.rotorPitchIndicator.SetRotorPitch(rotorPitch + dPitch)

	rotorRPM := h.rotorRPMIndicator.GetRotorRPM()
	if rotorRPM >= 100 {
		dRPM = -0.1
	} else if rotorRPM <= 80 {
		dRPM = 0.1
	}
	h.rotorRPMIndicator.SetRotorRPM(rotorRPM + dRPM)

	h.rotorPitchImg.Update(h.rotorPitchIndicator.GetImage())
	h.rotorRPMImg.Update(h.rotorRPMIndicator.GetImage())

	return nil
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

func (h *HUD) Draw(screen *ebiten.Image) {
	rotorPitchImg := &h.rotorPitchImg
	rotorRPMImg := &h.rotorRPMImg

	if !rotorPitchImg.NeedToDraw &&
		!rotorRPMImg.NeedToDraw {
		return
	}

	op := &ebiten.DrawImageOptions{CompositeMode: ebiten.CompositeModeCopy}
	op.GeoM.Translate(0, 2*rowHeight)
	rotorPitchImg.DrawOn(screen, op)

	op.GeoM.Translate(float64(rotorPitchImg.Size().X), 0)
	rotorRPMImg.DrawOn(screen, op)

	/*	ff := h.fontFace

		ff.DrawTextWithShadow(screen, "315", 0, 0, textColor)
		ff.DrawTextWithShadowCenter(screen, "108", 0, 0, textColor, ScreenWidth)
		ff.DrawTextWithShadowRight(screen, "12", 0, 0, textColor, ScreenWidth)
	*/
}

func (h *HUD) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}