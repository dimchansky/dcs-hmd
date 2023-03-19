package rotorpitch

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dimchansky/dcs-hmd/gui/indicator"
)

type IndicatorConfig struct {
	Width           int
	Height          int
	TickLength      int
	MinorTickLength int
	LineWidth       float64
	Color           color.NRGBA
	BorderColor     color.NRGBA
	Rect            image.Rectangle
}

const (
	minPitch  = 1
	maxPitch  = 15
	pitchStep = 1
)

func NewIndicator(cfg *IndicatorConfig) *Indicator {
	return &Indicator{
		impl: indicator.New(
			&indicator.Config{
				Width:               cfg.Width,
				Height:              cfg.Height,
				TickLength:          cfg.TickLength,
				MinorTickLength:     cfg.MinorTickLength,
				LineWidth:           cfg.LineWidth,
				Color:               cfg.Color,
				BorderColor:         cfg.BorderColor,
				Rect:                cfg.Rect,
				MinValue:            minPitch,
				MaxValue:            maxPitch,
				MinFixedWindowValue: minPitch,
				MaxFixedWindowValue: maxPitch,
				MinTickStep:         pitchStep,
				GetTickLength: func(rotorPitchValue int) float64 {
					tickLen := float64(cfg.MinorTickLength)
					if rotorPitchValue%2 != 0 {
						tickLen = float64(cfg.TickLength)
					}
					return tickLen
				},
				MinSafeValue:    nil,
				MaxAllowedValue: nil,
				MinLabelStep:    2,
				GetLabelOffset: func(rotorPitchValue int) float64 {
					return float64(cfg.TickLength)
				},
			},
		),
	}
}

type Indicator struct {
	impl *indicator.Indicator
}

func (i *Indicator) SetRotorPitch(rotorPitch float64) {
	i.impl.SetValue(rotorPitch)
}

func (i *Indicator) GetRotorPitch() (rotorPitch float64) {
	return i.impl.GetValue()
}

func (i *Indicator) GetImage() (img *ebiten.Image, isRedrawn bool) {
	img, isRedrawn = i.impl.GetImage()
	return
}
