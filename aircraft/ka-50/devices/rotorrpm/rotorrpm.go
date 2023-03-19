package rotorrpm

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

func NewIndicator(cfg *IndicatorConfig) *Indicator {
	// Maximum allowed rotor RPM – 98%
	// Minimum safe RPM in flight – 83%
	minSafeRPM := 83
	maxAllowedRPM := 98

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
				MinValue:            0,
				MaxValue:            110,
				MinFixedWindowValue: 80,
				MaxFixedWindowValue: 100,
				MinTickStep:         1,
				GetTickLength: func(rpmValue int) float64 {
					tickLen := float64(cfg.MinorTickLength / 2)
					if rpmValue%10 == 0 {
						tickLen = float64(cfg.TickLength)
					} else if rpmValue%5 == 0 {
						tickLen = float64(cfg.MinorTickLength)
					}
					return tickLen
				},
				MinSafeValue:    &minSafeRPM,
				MaxAllowedValue: &maxAllowedRPM,
				MinLabelStep:    10,
				GetLabelOffset: func(rpmValue int) float64 {
					return float64(cfg.TickLength)
				},
			},
		),
	}
}

type Indicator struct {
	impl *indicator.Indicator
}

func (i *Indicator) SetRotorRPM(rotorRPM float64) {
	i.impl.SetValue(rotorRPM)
}

func (i *Indicator) GetRotorRPM() (rotorRPM float64) {
	return i.impl.GetValue()
}

func (i *Indicator) GetImage() (img *ebiten.Image, isRedrawn bool) {
	img, isRedrawn = i.impl.GetImage()
	return
}
