package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Point struct {
	x, y float64
}

type Game struct {
	width, height, pow int
	p                  []Point
}

func (g *Game) Update() error {
	return nil
}

func gaussianElimination(n int, m [][]float64) []float64 {
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if m[j][i] != 0 {
				m[i], m[j] = m[j], m[i]
				break
			}
		}
		for k := i + 1; k < n; k++ {
			q := m[k][i] / m[i][i]
			for j := i; j < n+1; j++ {
				m[k][j] -= q * m[i][j]
			}
		}
	}

	ans := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		var sum float64
		for j := i + 1; j < n; j++ {
			sum += m[i][j] * ans[j]
		}
		ans[i] = (m[i][n] - sum) / m[i][i]
	}
	return ans
}

func approximation(n int, p []Point) []float64 {
	n += 1
	sX, sY := make([]float64, 2*n), make([]float64, n)
	sX[0] = float64(len(p))
	for _, v := range p {
		for i := 0; i < n; i++ {
			sY[i] += v.y * math.Pow(v.x, float64(i))
		}
		for i := 0; i < 2*n; i++ {
			sX[i] += math.Pow(v.x, float64(i))
		}
	}

	m := make([][]float64, n)
	for i := 0; i < n; i++ {
		m[i] = make([]float64, n+1)
		for j := 0; j < n+1; j++ {
			m[i][j] = sX[i+j]
		}
		m[i][n] = sY[i]
	}
	return gaussianElimination(n, m)
}

func solve(coef []float64, x float64) float64 {
	var ans float64
	for i := len(coef) - 1; i >= 0; i-- {
		ans += coef[i] * math.Pow(x, float64(i))
	}
	return ans
}

func (g *Game) Draw(screen *ebiten.Image) {
	coef := approximation(g.pow, g.p)
	for i := 0.0; i < float64(g.width); i += 0.1 {
		vector.DrawFilledCircle(screen, float32(i), float32(solve(coef, i)), 1, color.RGBA{255, 0, 255, 255}, true)
	}
	for _, i := range g.p {
		vector.DrawFilledCircle(screen, float32(i.x), float32(i.y), 2, color.RGBA{255, 255, 0, 0}, true)
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func NewGame(width, height, pow int, p []Point) *Game {
	return &Game{
		width:  width,
		height: height,
		pow:    pow,
		p:      p,
	}
}

func main() {
	var pow, np int
	fmt.Println("Enter power and number of points:")
	fmt.Scan(&pow, &np)
	var p []Point
	for i := 0; i < np; i++ {
		p = append(p, Point{rand.Float64() * 800, rand.Float64() * 600})
	}
	g := NewGame(screenWidth, screenHeight, pow, p)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
