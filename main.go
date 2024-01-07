package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

const (
	screenWidth, screenHeight                    = 1800, 1080
	plotMinX, plotMaxX, plotMinY, plotMaxY       = 0, 10, 0, 100 // Min and Max data values along both axis
	pointMinYOffset, pointMaxYOffset, pointCount = -20, 20, 10
)

func f(x float64) float64 { return x * x } // Function to spawn points along

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Approximator")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

type game struct {
	points  plotter.XYs
	NewPlot func() *plot.Plot
	plot    *plot.Plot
}

func NewGame() *game {
	return &game{
		GetRandPoints(f, pointMinYOffset, pointMaxYOffset, pointCount),
		func() *plot.Plot {
			p := plot.New()
			p.X.Min = plotMinX
			p.X.Max = plotMaxX
			p.Y.Min = plotMinY
			p.Y.Max = plotMaxY

			p.BackgroundColor = color.Black
			return p
		},
		nil,
	}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }
func (g *game) Update() error {
	g.plot = g.NewPlot()

	s, err := plotter.NewScatter(g.points)
	if err != nil {
		return err
	}
	s.Color = color.RGBA{255, 0, 0, 255}
	g.plot.Add(s)

	original := plotter.NewFunction(f)
	original.Color = color.RGBA{100, 100, 100, 255}
	original.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	g.plot.Add(original)

	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	DrawPlot(screen, g.plot)
}

// Returns random points along f with random Yoffset
func GetRandPoints(f func(float64) float64, minYoffset float64, maxYoffset float64, pointCount uint) (p plotter.XYs) {
	// 1. Getting random argument value X
	// 2. Getting function value(Y)
	// 3. Applying offset to Y
	for i := uint(0); i < pointCount; i++ {
		x := plotMinX + rand.Float64()*(plotMaxX-plotMinX) // Random argument within visible range
		yOffset := minYoffset + rand.Float64()*(maxYoffset-minYoffset)
		p = append(p, plotter.XY{x, f(x) + yOffset})
	}
	return
}

func DrawPlot(screen *ebiten.Image, p *plot.Plot) {
	// https://github.com/gonum/plot/wiki/Drawing-to-an-Image-or-Writer:-How-to-save-a-plot-to-an-image.Image-or-an-io.Writer,-not-a-file.
	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	c := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(c))

	screen.DrawImage(ebiten.NewImageFromImage(c.Image()), &ebiten.DrawImageOptions{})
}

// Returns approximating linear function
func Approximate(points plotter.XYs) func(float64) float64 {
	var sx, sxx, sy, sxy float64
	n := float64(points.Len())
	for _, p := range points {
		sx, sxx, sy, sxy = sx+p.X, sxx+p.X*p.X, sy+p.Y, sxy+p.X*p.Y
	}
	a := (n*sxy - sx*sy) / (n*sxx - sx*sx)
	return func(x float64) float64 {
		return a*x + (sy-a*sx)/n
	}
}

// The task:
// 1. Generate random points along some function
// Game struct contains function to draw points along, Max and min Y offset, point count, power of polynomial to approximate with
// 2. Calculate coeffecients of the approximating polynomial
// 3. Solve system of equations using Gaussian Elimination
