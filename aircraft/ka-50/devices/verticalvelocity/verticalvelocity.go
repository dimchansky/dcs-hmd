package verticalvelocity

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
				MinValue:            -30,
				MaxValue:            30,
				MinFixedWindowValue: -8,
				MaxFixedWindowValue: 8,
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
				MinSafeValue:    nil,
				MaxAllowedValue: nil,
				MinLabelStep:    5,
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

func (i *Indicator) SetVerticalVelocity(verticalVelocity float64) {
	i.impl.SetValue(verticalVelocity)
}

func (i *Indicator) GetVerticalVelocity() (verticalVelocity float64) {
	return i.impl.GetValue()
}

func (i *Indicator) GetImage() (img *ebiten.Image, isRedrawn bool) {
	img, isRedrawn = i.impl.GetImage()
	return
}
