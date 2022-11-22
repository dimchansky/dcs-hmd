package ka50rotorpitch

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

	// draw indicator gauge
	dc := gg.NewContext(width, height)

	verticalLineX := float64((maxPoint.X-minPoint.X)/2 + minPoint.X)
	yTop := float64(minPoint.Y)
	yBottom := float64(maxPoint.Y)

	dc.DrawLine(verticalLineX, yTop, verticalLineX, yBottom)

	const (
		minPitch  = 1
		maxPitch  = 15
		pitchStep = 1
	)

	rotorPitchToY := &utils.IntervalTransformer{
		Interval1: utils.Interval{
			Start: minPitch,
			End:   maxPitch,
		},
		Interval2: utils.Interval{
			Start: yBottom,
			End:   yTop,
		},
	}

	for pitch := minPitch; pitch <= maxPitch; pitch += pitchStep {
		x := verticalLineX - float64(cfg.MinorTickLength)
		if pitch%2 != 0 {
			x = verticalLineX - float64(cfg.TickLength)
		}

		y := rotorPitchToY.TransformForward(float64(pitch))
		dc.DrawLine(x, y, verticalLineX, y)
	}

	dc.SetColor(cfg.BorderColor)
	dc.SetLineWidth(cfg.LineWidth * 3)
	dc.StrokePreserve()

	dc.SetColor(cfg.Color)
	dc.SetLineWidth(cfg.LineWidth)
	dc.Stroke()

	// draw labels
	for pitch := minPitch; pitch <= maxPitch; pitch += 2 {
		label := strconv.Itoa(pitch)
		x := verticalLineX - float64(cfg.TickLength)
		y := rotorPitchToY.TransformForward(float64(pitch))

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
		verticalLineX: verticalLineX,
		rotorPitchToY: rotorPitchToY,
	}

	currentRotorPitch := rotorPitchToY.Interval1.Start
	i.SetRotorPitch(currentRotorPitch)
	i.redrawFinalImage(currentRotorPitch)

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
	rotorPitchToY *utils.IntervalTransformer

	// drawn state
	drawnRotorPitch float64

	// thread-safe
	rotorPitchToDraw float64
}

func (i *Indicator) SetRotorPitch(rotorPitch float64) {
	rotorPitch = i.rotorPitchToY.Interval1.Sat(rotorPitch)

	m := &i.rwMutex
	m.Lock()
	i.rotorPitchToDraw = rotorPitch
	m.Unlock()
}

func (i *Indicator) GetRotorPitch() (rotorPitch float64) {
	m := &i.rwMutex
	m.RLock()
	rotorPitch = i.rotorPitchToDraw
	m.RUnlock()

	return
}

func (i *Indicator) GetImage() (img *ebiten.Image, isRedrawn bool) {
	// optimization: redraw the final image only if the value has changed
	if rotorPitchToDraw := i.GetRotorPitch(); rotorPitchToDraw != i.drawnRotorPitch {
		i.redrawFinalImage(rotorPitchToDraw)

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
	rotorPitchY := i.rotorPitchToY.TransformForward(rotorPitch)
	op.GeoM.Translate(i.verticalLineX, rotorPitchY)
	op.GeoM.Translate(-i.handPoint.X, -i.handPoint.Y)
	finalImg.DrawImage(i.handImg, op)

	// update the value for which the final image is rendered
	i.drawnRotorPitch = rotorPitch
}
