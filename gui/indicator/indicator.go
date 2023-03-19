package indicator

import (
	"image"
	"image/color"
	"math"
	"strconv"
	"sync"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/dimchansky/dcs-hmd/utils"
)

type Config struct {
	Width               int
	Height              int
	TickLength          int
	MinorTickLength     int
	LineWidth           float64
	Color               color.NRGBA
	BorderColor         color.NRGBA
	Rect                image.Rectangle
	MinValue            int
	MaxValue            int
	MinFixedWindowValue int
	MaxFixedWindowValue int
	MinTickStep         int
	GetTickLength       func(value int) float64
	MinSafeValue        *int
	MaxAllowedValue     *int
	MinLabelStep        int
	GetLabelOffset      func(value int) float64
}

func New(cfg *Config) *Indicator {
	width := cfg.Width
	height := cfg.Height

	// coordinates of the upper left corner and the lower right corner of the window in which the indicator is to be displayed
	maxPoint := cfg.Rect.Max
	minPoint := cfg.Rect.Min
	fixedWindowHeight := cfg.Rect.Dy()

	// Only a small part of the indicator gauge will be displayed in the window.
	// Values from MinFixedWindowValue to MaxFixedWindowValue will be displayed in a window with a fixed indicator gauge,
	// in the range of these values will move the hand. If the values are out of these limits, the indicator gauge will
	// move and the hand will freeze in the upper or lower position, depending on where the values have gone.

	// Calculate the necessary height of the inner window, which would fit the entire range of values.
	// (MaxFixedWindowValue - MinFixedWindowValue) ~ fixedWindowScreenHeight pixels
	// (MaxValue - MinValue)                       ~        ???        pixels
	fullRangeGaugeHeight := int(math.Ceil(float64(cfg.MaxValue-cfg.MinValue) * float64(fixedWindowHeight) / float64(cfg.MaxFixedWindowValue-cfg.MinFixedWindowValue)))

	// Calculate by how much the total height of the indicator strip should increase
	gaugeHeight := height + (fullRangeGaugeHeight - fixedWindowHeight)

	// draw indicator gauge
	dc := gg.NewContext(width, gaugeHeight)

	verticalLineX := utils.Transform(1,
		&utils.Interval{Start: 0, End: 2},
		&utils.Interval{Start: float64(minPoint.X), End: float64(maxPoint.X - 1)},
	)

	yTop := float64(minPoint.Y)
	yBottom := float64(minPoint.Y) + float64(fullRangeGaugeHeight-1)
	dc.DrawLine(verticalLineX, yTop, verticalLineX, yBottom)

	valueToScreenY := &utils.IntervalTransformer{
		IntervalFrom: utils.Interval{
			Start: float64(cfg.MinValue),
			End:   float64(cfg.MaxValue),
		},
		IntervalTo: utils.Interval{
			Start: yBottom,
			End:   yTop,
		},
	}

	for value := cfg.MinValue; value <= cfg.MaxValue; value += cfg.MinTickStep {
		x1 := verticalLineX - cfg.GetTickLength(value)
		y := valueToScreenY.TransformForward(float64(value))

		dc.DrawLine(x1, y, verticalLineX, y)
	}

	x3 := verticalLineX + float64(cfg.TickLength)
	if minSafeValue := cfg.MinSafeValue; minSafeValue != nil {
		minSafeValueScreenY := valueToScreenY.TransformForward(float64(*minSafeValue))
		dc.MoveTo(verticalLineX, minSafeValueScreenY)
		dc.LineTo(x3, minSafeValueScreenY)
		dc.LineTo(x3, minSafeValueScreenY-float64(cfg.TickLength))
	}

	if maxAllowedValue := cfg.MaxAllowedValue; maxAllowedValue != nil {
		maxAllowedValueScreenY := valueToScreenY.TransformForward(float64(*maxAllowedValue))
		dc.MoveTo(verticalLineX, maxAllowedValueScreenY)
		dc.LineTo(x3, maxAllowedValueScreenY)
		dc.LineTo(x3, maxAllowedValueScreenY+float64(cfg.TickLength))
	}

	dc.SetColor(cfg.BorderColor)
	dc.SetLineWidth(cfg.LineWidth * 3)
	dc.StrokePreserve()

	dc.SetColor(cfg.Color)
	dc.SetLineWidth(cfg.LineWidth)
	dc.Stroke()

	// draw labels
	for labelValue := cfg.MinValue; labelValue <= cfg.MaxValue; labelValue += cfg.MinLabelStep {
		label := strconv.Itoa(labelValue)
		x := verticalLineX - cfg.GetLabelOffset(labelValue)
		y := valueToScreenY.TransformForward(float64(labelValue))

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
		finalImg:                   ebiten.NewImage(width, height),
		gaugeImg:                   gaugeImg,
		handImg:                    handImg,
		handPoint:                  handPoint,
		verticalLineX:              verticalLineX,
		maxFixedWindowValue:        float64(cfg.MaxFixedWindowValue),
		minFixedWindowValue:        float64(cfg.MinFixedWindowValue),
		valueToScreenY:             valueToScreenY,
		maxValueScreenY:            valueToScreenY.TransformForward(float64(cfg.MaxValue)),
		maxFixedWindowValueScreenY: valueToScreenY.TransformForward(float64(cfg.MaxFixedWindowValue)),
		fixedWindowScreenHeight:    valueToScreenY.TransformForward(float64(cfg.MinFixedWindowValue)) - valueToScreenY.TransformForward(float64(cfg.MaxFixedWindowValue)),
	}

	currentValue := valueToScreenY.IntervalFrom.Start
	i.SetValue(currentValue)
	i.redrawFinalImage(currentValue)

	return i
}

type Indicator struct {
	rwMutex sync.RWMutex

	// images
	finalImg *ebiten.Image
	gaugeImg *ebiten.Image
	handImg  *ebiten.Image

	// image transformation variables
	handPoint      gg.Point
	verticalLineX  float64
	valueToScreenY *utils.IntervalTransformer

	// drawn state
	drawnValue float64

	// thread-safe
	maxFixedWindowValue        float64
	minFixedWindowValue        float64
	valueToDraw                float64
	maxValueScreenY            float64
	maxFixedWindowValueScreenY float64
	fixedWindowScreenHeight    float64
}

func (i *Indicator) SetValue(value float64) {
	value = i.valueToScreenY.IntervalFrom.Sat(value)

	m := &i.rwMutex
	m.Lock()
	i.valueToDraw = value
	m.Unlock()
}

func (i *Indicator) GetValue() (value float64) {
	m := &i.rwMutex
	m.RLock()
	value = i.valueToDraw
	m.RUnlock()

	return
}

func (i *Indicator) GetImage() (img *ebiten.Image, isRedrawn bool) {
	// optimization: redraw the final image only if the value has changed
	if valueToDraw := i.GetValue(); valueToDraw != i.drawnValue {
		i.redrawFinalImage(valueToDraw)

		isRedrawn = true
	}

	img = i.finalImg

	return
}

func (i *Indicator) redrawFinalImage(value float64) {
	finalImg := i.finalImg
	finalImg.Clear()

	valueToScreenY := i.valueToScreenY
	valueScreenY := valueToScreenY.TransformForward(value)

	// draw gauge
	op := &ebiten.DrawImageOptions{}
	var gaugeXTranslate float64
	if value > i.maxFixedWindowValue {
		gaugeXTranslate = i.maxValueScreenY - valueScreenY
	} else if value < i.minFixedWindowValue {
		gaugeXTranslate = i.maxValueScreenY - valueScreenY + i.fixedWindowScreenHeight
	} else {
		// fixed window case
		gaugeXTranslate = i.maxValueScreenY - i.maxFixedWindowValueScreenY
	}
	op.GeoM.Translate(0, gaugeXTranslate)
	finalImg.DrawImage(i.gaugeImg, op)

	// draw hand
	op.GeoM.Translate(i.verticalLineX, valueScreenY)
	op.GeoM.Translate(-i.handPoint.X, -i.handPoint.Y)
	finalImg.DrawImage(i.handImg, op)

	// update the value for which the final image is rendered
	i.drawnValue = value
}
