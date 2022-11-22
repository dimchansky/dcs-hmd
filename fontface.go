package dcshmd

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func NewFontFace(size int, dpi int) (*FontFace, error) {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PressStart2P.ttf font: %w", err)
	}

	ff, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     float64(dpi),
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new font face: %w", err)
	}

	return &FontFace{
		fontBaseSize: size,
		Face:         ff,
	}, nil
}

type FontFace struct {
	fontBaseSize int
	font.Face
}

func (h *FontFace) TextWidth(str string) int {
	maxW := 0
	for _, line := range strings.Split(str, "\n") {
		b, _ := font.BoundString(h.Face, line)
		w := (b.Max.X - b.Min.X).Ceil()
		if maxW < w {
			maxW = w
		}
	}
	return maxW
}

func (h *FontFace) DrawText(rt *ebiten.Image, str string, x, y int, clr color.Color) {
	offsetY := h.fontBaseSize
	for _, line := range strings.Split(str, "\n") {
		y += offsetY
		text.Draw(rt, line, h.Face, x, y, clr)
	}
}

func (h *FontFace) DrawTextWithShadow(rt *ebiten.Image, str string, x, y int, clr color.Color) {
	offsetY := h.fontBaseSize
	for _, line := range strings.Split(str, "\n") {
		y += offsetY
		text.Draw(rt, line, h.Face, x+1, y+1, shadowColor)
		text.Draw(rt, line, h.Face, x, y, clr)
	}
}

func (h *FontFace) DrawTextWithShadowCenter(rt *ebiten.Image, str string, x, y int, clr color.Color, width int) {
	w := h.TextWidth(str)
	x += (width - w) / 2
	h.DrawTextWithShadow(rt, str, x, y, clr)
}

func (h *FontFace) DrawTextWithShadowRight(rt *ebiten.Image, str string, x, y int, clr color.Color, width int) {
	w := h.TextWidth(str)
	x += width - w
	h.DrawTextWithShadow(rt, str, x, y, clr)
}
