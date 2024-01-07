package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type coord struct {
	x float64
	y float64
}

// Game implements ebiten.Game interface.
type Game struct {
	points []coord
	n      int
}

const (
	screenWidth  = 640
	screenHeight = 480
)

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.points = append(g.points, coord{float64(x), float64(y)})
	}
	// Write your game's logical update.
	return nil
}

func f(c []float64, x float64) (res float64) {
	for i := len(c) - 1; i >= 0; i-- {
		res += c[i] * math.Pow(x, float64(i))
	}
	return
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	var c []float64
	if len(g.points) > 1 {
		sumX := make([]float64, g.n+2)
		sumX[0] = float64(len(g.points))
		var sumx float64
		for _, v := range g.points {
			sumx += v.x
		}
		sumX[1] = sumx
		for i := 2; i <= g.n+1; i++ {
			sumX[i] = sumX[i-1] * sumx
		}
		m := make([][]float64, g.n)
		for i := range m {
			m[i] = make([]float64, g.n+1)
			for j := range g.points {
				m[i][g.n] += math.Pow(g.points[j].x, float64(g.n-i-1)) * g.points[j].y
			}
		}
		a := g.n + 1
		for i := 0; i < g.n; i++ {
			b := a - i
			for j := 0; j < g.n; j++ {
				fmt.Println(b)
				m[i][j] = sumX[b]
				b--
			}
		}
		printMatrix(m)
		fmt.Println(len(m))
		Eliminate(m)
		c = solve(m)
		fmt.Println(c)
	}
	for x := 0; x < screenWidth; x++ {
		y := f(c, float64(x))
		screen.Set(x, int(y)+screenHeight/2, color.White)
	}

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{[]coord{}, 3}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Aproximation")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
