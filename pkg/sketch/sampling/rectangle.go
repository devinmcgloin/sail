package sampling

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/devinmcgloin/clr/clr"
	"github.com/devinmcgloin/sail/pkg/canvas"
	"github.com/devinmcgloin/sail/pkg/fill"
	"github.com/devinmcgloin/sail/pkg/shapes"
	"github.com/fogleman/gg"
)

type UniformRectangleDot struct{}

func (c UniformRectangleDot) Dimensions() (int, int) {
	return 1400, 900
}

func (c UniformRectangleDot) Draw(context *gg.Context, r *rand.Rand) {
	rows := 1 + math.Floor(r.Float64()*9)
	margin := r.Float64() * 0.10
	startHue := r.Intn(365)
	endHue := r.Intn(365)
	gaps := r.Float64() * 0.5
	fmt.Printf("\tGapFactor: %f\n", gaps)

	filler := fill.NewUniformFiller(8000, r)

	for i := 0.0; i < rows; i++ {
		rect := shapes.Rectangle{
			A: shapes.Point{
				X: canvas.W(context, rectangePositioning(margin, i, rows)),
				Y: canvas.H(context, margin),
			},
			B: shapes.Point{
				X: canvas.W(context, rectangePositioning(margin, 1+i, rows)),
				Y: canvas.H(context, 1.0-margin),
			},
		}

		rect.ShrinkHorizontally(gaps)

		hsv := clr.HSV{H: startHue + int(float64(endHue)/(rows-1)*i), S: 60, V: int(20 + i*(80.0/(rows-1)))}
		r, g, b := hsv.RGB()
		context.SetRGB(float64(r), float64(g), float64(b))

		filler.DotFill(context, rect)
	}
}

type RadialRectangleDot struct {
}

func (c RadialRectangleDot) Dimensions() (int, int) {
	return 1400, 900
}

func (c RadialRectangleDot) Draw(context *gg.Context, r *rand.Rand) {
	rows := 1 + math.Floor(r.Float64()*15)
	margin := r.Float64() * 0.10
	hue := r.Intn(365)
	gaps := r.Float64()

	fmt.Printf("\tGapFactor: %f\n", gaps)

	filler := fill.NewRadialFiller(8000, r)

	for i := 0.0; i < rows; i++ {
		rect := shapes.Rectangle{
			A: shapes.Point{
				X: canvas.W(context, rectangePositioning(margin, i, rows)),
				Y: canvas.H(context, margin),
			},
			B: shapes.Point{
				X: canvas.W(context, rectangePositioning(margin, 1+i, rows)),
				Y: canvas.H(context, 1.0-margin),
			},
		}
		rect.ShrinkHorizontally(gaps)

		r, g, b := clr.HSV{H: hue, S: int(i * 7), V: 70}.RGB()
		context.SetRGB(float64(r), float64(g), float64(b))

		filler.Fill(context, rect)
	}
}

func rectangePositioning(offset, index, rectangeCount float64) float64 {
	avaliableSpace := 1.0 - offset*2
	return (offset + index*(avaliableSpace/rectangeCount))
}
