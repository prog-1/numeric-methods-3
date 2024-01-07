package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
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
	power                                        = 1 // The power of polynomial to appriximate with
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
	points                plotter.XYs
	approximatingFunction func(float64) float64
	NewPlot               func() *plot.Plot
	plot                  *plot.Plot
}

func NewGame() *game {
	points := GetRandPoints(f, pointMinYOffset, pointMaxYOffset, pointCount)
	return &game{
		points,
		Approximate(points),
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

	approximating := plotter.NewFunction(g.approximatingFunction)
	approximating.Color = color.RGBA{255, 255, 255, 255}
	g.plot.Add(approximating)

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
	// 1. Compose the matrix
	// 1.1. Get the size of a matrix
	// 1.2.
	// 2. Solve the matrix using Gaussian Elimination Check

	// Composing the matrix
	A := make([][]float64, power+1) // Rows
	for i := 0; i < len(A); i++ {
		A[i] = make([]float64, power+2) // Columns
		// Matrix vectors
		for j := 0; j < len(A[i])-1; j++ {
			for k := 0; k < points.Len(); k++ {
				A[i][j] += math.Pow(points[k].X, float64(power)*2-float64(j+i))
			}
		}
		// Output vector
		for k := 0; k < points.Len(); k++ {
			A[i][power+1] = points[k].Y * math.Pow(points[k].X, float64(power-k))
		}
	}

	// n := points.Len()
	// var sxx, sx, sy, sxy float64
	// for i := 0; i < n; i++ {
	// 	x, y := points[i].X, points[i].Y
	// 	sxx += x * x
	// 	sx += x
	// 	sy += y
	// 	sxy += x * y
	// }

	// Calculating coefecients
	coefs := gauss(A)
	l := len(coefs)

	// Printing approximating function
	fmt.Print("Approximating function: ")
	for i, c := range coefs {
		fmt.Printf("%vx^%v + ", c, l-i-1.)
	}
	fmt.Println()

	// Returning approximating function
	return func(x float64) (res float64) {
		for i, c := range coefs {
			res += c * math.Pow(x, float64(l-i-1))
		}
		return res
	}
}

// Solves system of linear equations written in the matrix form in the row-major order
// and with output vector put to the right side of matrix
// Returns coeffecients from top to bottom(x, y, z ...)// Solves system of linear equations written in the matrix form
// https://github.com/34thSchool/numeric-methods-2
func gauss(A [][]float64) []float64 {
	n := len(A)
	// 1. Finding a row with in eliminated coefficient
	for i := 0; i < n; i++ {
		if A[i][0] != 0 {
			A[0], A[i] = A[i], A[0]
			break
		}
	}
	// 2. Eliminating
	for j := 0; j < n-1; j++ { // Columns
		for i := j + 1; i < n; i++ { // Rows
			c := A[i][j] / A[j][j]    // Subtrahend column multiplier
			for k := j; k <= n; k++ { // Columns
				A[i][k] -= c * A[j][k] // Elimination
			}
		}
	}
	// 3. Solving
	res := make([]float64, len(A))
	for i := n - 1; i >= 0; i-- { // Rows
		res[i] = A[i][n]
		for j := i + 1; j < n; j++ {
			res[i] -= A[i][j] * res[j]
		}
		res[i] /= A[i][i]
	}
	return res
}
