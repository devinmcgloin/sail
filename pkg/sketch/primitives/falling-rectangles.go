package primitives

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/devinmcgloin/sail/pkg/slog"
	"github.com/fogleman/gg"
)

//FallingRectangles defines the sketch type
type FallingRectangles struct{}

// Dimensions returns the size of the sketch in integers
func (fr FallingRectangles) Dimensions() (int, int) {
	return int(fr.Width()), int(fr.Height())
}

// Width helper for finding the width of the canvas
func (fr FallingRectangles) Width() float64 {
	return 12000.0
}

// Height helper for finding the hieght of the canvas
func (fr FallingRectangles) Height() float64 {
	return 19500.0
}

// Draw is the primary rendering method
func (fr FallingRectangles) Draw(context *gg.Context, rand *rand.Rand) {
	margin := rand.Float64()*(fr.Width()*0.15) + fr.Width()*0.04

	avaliableSpace := fr.Width() - margin*2
	sizeFactor := math.Floor(rand.Float64()*25 + 5)
	boxSize := avaliableSpace / sizeFactor
	halfBox := boxSize / 2

	noiseFactor := rand.Float64() * 4

	context.SetLineWidth(fr.Width() * 0.001)
	context.SetColor(color.Black)

	slog.InfoValues("margin", margin, "avaliableSpace", avaliableSpace, "sizeFactor", sizeFactor, "noiseFactor", noiseFactor, "boxSize", boxSize)

	for x := margin + halfBox; x < fr.Width()-margin; x += boxSize {
		rowIndex := 0.0
		for y := margin + halfBox; y < fr.Height()-margin; y += boxSize {
			sectionOffset := rowIndex * boxSize
			normalizedNoise := (sectionOffset / (fr.Height() - margin))
			if normalizedNoise > 1 {
				slog.InfoValues("sectionOffset", sectionOffset, "normalizedNoise", normalizedNoise, "fr.Height() - margin", fr.Height()-margin)
			}
			adjustedNoise := normalizedNoise * noiseFactor
			rotate := adjustedNoise * rand.NormFloat64()
			context.Push()
			context.RotateAbout(rotate, x, y)
			translationFactor := adjustedNoise * noiseFactor * fr.Height() * 0.01
			context.Translate(rand.Float64()*translationFactor, rand.Float64()*translationFactor)
			context.DrawRectangle(x-halfBox, y-halfBox, boxSize, boxSize)
			context.Stroke()
			context.Pop()
			rowIndex++
		}
	}
}

func drawAbsRect(dc *gg.Context, x1, y1, x2, y2 float64) {
	dc.NewSubPath()
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x1, y2)
	dc.ClosePath()
}
