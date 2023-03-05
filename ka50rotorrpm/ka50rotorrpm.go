package ka50rotorrpm

import (
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dimchansky/dcs-hmd/utils"
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
	width := cfg.Width
	height := cfg.Height
	maxPoint := cfg.Rect.Max
	minPoint := cfg.Rect.Min

	dc := gg.NewContext(width, height)

	x2 := float64((maxPoint.X-minPoint.X)/2 + minPoint.X)
	yTop := float64(minPoint.Y)
	yBottom := float64(maxPoint.Y)

	dc.DrawLine(x2, yTop, x2, yBottom)

	// Maximum allowed rotor RPM – 98%
	// Minimum safe RPM in flight – 83%
	const (
		minRPM        = 80
		maxRPM        = 100
		minSafeRPM    = 83
		maxAllowedRPM = 98
		rpmStep       = 1
	)

	rotorRPMToY := &utils.IntervalTransformer{
		IntervalFrom: utils.Interval{
			Start: minRPM,
			End:   maxRPM,
		},
		IntervalTo: utils.Interval{
			Start: yBottom,
			End:   yTop,
		},
	}

	for rpm := minRPM; rpm <= maxRPM; rpm += rpmStep {
		x1 := x2 - float64(cfg.MinorTickLength/2)
		if rpm%10 == 0 {
			x1 = x2 - float64(cfg.TickLength)
		} else if rpm%5 == 0 {
			x1 = x2 - float64(cfg.MinorTickLength)
		}

		y := rotorRPMToY.TransformForward(float64(rpm))
		dc.DrawLine(x1, y, x2, y)
	}

	x3 := x2 + float64(cfg.TickLength)
	rpm83Y := rotorRPMToY.TransformForward(minSafeRPM)
	dc.MoveTo(x2, rpm83Y)
	dc.LineTo(x3, rpm83Y)
	dc.LineTo(x3, rpm83Y+float64(cfg.TickLength))

	rpm98Y := rotorRPMToY.TransformForward(maxAllowedRPM)
	dc.MoveTo(x2, rpm98Y)
	dc.LineTo(x3, rpm98Y)
	dc.LineTo(x3, rpm98Y-float64(cfg.TickLength))

	dc.SetColor(cfg.BorderColor)
	dc.SetLineWidth(cfg.LineWidth * 3)
	dc.StrokePreserve()

	dc.SetColor(cfg.Color)
	dc.SetLineWidth(cfg.LineWidth)
	dc.Stroke()

	// draw labels
	for rpm := minRPM; rpm <= maxRPM; rpm += 10 {
		label := strconv.Itoa(rpm)
		x := x2 - float64(cfg.TickLength)
		y := rotorRPMToY.TransformForward(float64(rpm))

		dc.SetColor(cfg.BorderColor)

		const (
			n  = 3 // "stroke" size
			ax = 0.3
			ay = 0.4
		)

		for dy := -n; dy <= n; dy++ {
			for dx := -n; dx <= n; dx++ {
				if dx*dx+dy*dy >= n*n {
					// give it rounded corners
					continue
				}

				dc.DrawStringAnchored(label, x+float64(dx), y+float64(dy), ax, ay)
			}
		}
		dc.SetColor(cfg.Color)
		dc.DrawStringAnchored(label, x, y, ax, ay)
	}

	gaugeImg := ebiten.NewImageFromImage(dc.Image())

	// draw hand
	const (
		handSpan = 3
	)

	dc = gg.NewContext(cfg.TickLength+2*handSpan, cfg.TickLength+2*handSpan)

	dc.MoveTo(float64(cfg.TickLength+handSpan), handSpan)

	handPoint := gg.Point{
		X: handSpan,
		Y: handSpan + float64(cfg.TickLength)/2.0,
	}
	dc.LineTo(handPoint.X, handPoint.Y)
	dc.LineTo(float64(cfg.TickLength+handSpan), float64(cfg.TickLength+handSpan))
	dc.LineTo(float64(cfg.TickLength+handSpan), handSpan)

	dc.SetColor(cfg.BorderColor)
	dc.SetLineWidth(cfg.LineWidth * 3)
	dc.StrokePreserve()

	dc.SetColor(cfg.Color)
	dc.SetLineWidth(cfg.LineWidth)
	dc.Stroke()

	handImg := ebiten.NewImageFromImage(dc.Image())

	i := &Indicator{
		finalImg:      ebiten.NewImage(width, height),
		gaugeImg:      gaugeImg,
		handImg:       handImg,
		handPoint:     handPoint,
		verticalLineX: x2,
		rotorRPMToY:   rotorRPMToY,
	}

	currentRotorRPM := rotorRPMToY.IntervalFrom.Start
	i.SetRotorRPM(currentRotorRPM)
	i.redrawFinalImage(currentRotorRPM)

	return i
}

type Indicator struct {
	rwMutex sync.RWMutex

	// images
	finalImg *ebiten.Image
	gaugeImg *ebiten.Image
	handImg  *ebiten.Image

	// image transformation variables
	handPoint     gg.Point
	verticalLineX float64
	rotorRPMToY   *utils.IntervalTransformer

	// drawn state
	drawnRotorRPM float64

	// thread-safe
	rotorRPMToDraw float64
}

func (i *Indicator) SetRotorRPM(rotorRPM float64) {
	rotorRPM = i.rotorRPMToY.IntervalFrom.Sat(rotorRPM)

	m := &i.rwMutex
	m.Lock()
	i.rotorRPMToDraw = rotorRPM
	m.Unlock()
}

func (i *Indicator) GetRotorRPM() (rotorRPM float64) {
	m := &i.rwMutex
	m.RLock()
	rotorRPM = i.rotorRPMToDraw
	m.RUnlock()

	return
}

func (i *Indicator) GetImage() (img *ebiten.Image, isRedrawn bool) {
	// optimization: redraw the final image only if the value has changed
	if rotorRPMToDraw := i.GetRotorRPM(); rotorRPMToDraw != i.drawnRotorRPM {
		i.redrawFinalImage(rotorRPMToDraw)

		isRedrawn = true
	}

	img = i.finalImg

	return
}

func (i *Indicator) redrawFinalImage(rotorPitch float64) {
	finalImg := i.finalImg
	finalImg.Clear()

	// draw gauge
	op := &ebiten.DrawImageOptions{}
	finalImg.DrawImage(i.gaugeImg, op)

	// draw hand
	rotorPitchY := i.rotorRPMToY.TransformForward(rotorPitch)
	op.GeoM.Translate(i.verticalLineX, rotorPitchY)
	op.GeoM.Translate(-i.handPoint.X, -i.handPoint.Y)
	finalImg.DrawImage(i.handImg, op)

	// update the value for which the final image is rendered
	i.drawnRotorRPM = rotorPitch
}
