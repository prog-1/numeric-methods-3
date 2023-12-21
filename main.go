package main

import (
	"fmt"
    "math"
    "math/rand"
    "image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 960
	screenHeight = 720
)

type Point struct {
	x, y float64
}

type Game struct {
	width, height int
	points        []Point
	pow 		  int
}

func NewGame(width, height int, p []Point, power int) *Game {
	return &Game{
		width:  width,
		height: height,
		points: p,
		pow:	power,
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    for _, v := range g.points {
		vector.DrawFilledCircle(screen, float32(v.x), float32(v.y), 3, color.RGBA{255, 0, 0, 255}, true)
	}
    c := approximation(g.pow, g.points)
 	for i := 0.0; i < float64(g.width); i += 0.05 {
       	vector.DrawFilledCircle(screen, float32(i), float32(f(c, i)), 3, color.RGBA{255, 0, 0, 255}, true)
    }
}

func approximation(n int, p []Point) []float64 { // biggest ^ + 1
    sumsX, sumsY := make([]float64, (n-1)*2+1), make([]float64, n)
	sumsX[0] = float64(n)
    for _, v := range p {
        for i := 1; i < (n-1)*2+1; i++ {
            sumsX[i] += math.Pow(v.x, float64(i))
        }
        for i := 0; i < n; i++ {
            sumsY[i] += v.y*math.Pow(v.x, float64(i))
        }
    }
    m := make([][]float64, n)
    for i := 0; i < n; i++ {
        m[i] = make([]float64, n+1)
        for j := 0 ; j < n+1; j++ {
            if j != n {
                m[i][j] = sumsX[len(sumsX)-1-j-i]
            } else {
                m[i][j] = sumsY[len(sumsY)-1-i]
            }
        }
    }
    return gauss(n, m)
}

func gauss(n int, m [][]float64) (res []float64) {
	for i := 0; i < n; i++ {
		if m[i][i] == 0 {
			for l := i + 1; l < n; l++ {
				if m[l][l] != 0 {
					m[i], m[l] = m[l], m[i]
				}
			}
		}
		for j := i + 1; j < n; j++ {
			coef := m[j][i] / m[i][i]
			for l := i; l < n+1; l++ {
				m[j][l] -= m[i][l] * coef
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		for l := n - 1; l > i; l-- {
			m[i][n] = m[i][n] - m[i][l]
		}
		m[i][n] /= m[i][i]
		res = append(res, m[i][n])
		for j := i - 1; j >= 0; j-- {
			m[j][i] = m[j][i] * m[i][n]
		}
	}
	return
}

func f(c []float64, x float64) (res float64) {
    for i := len(c)-1; i >= 0; i-- {
        res += c[i]*math.Pow(x, float64(i))
    }
    return
}

func main() {
	var p []Point
    for i := 0; i < 10; i++ {
        p = append(p, Point{rand.Float64()*960, rand.Float64()*720})
    }
	var pow int
	fmt.Println("Enter power of polynomial: ")
	fmt.Scan(&pow)
	g := NewGame(screenWidth, screenHeight, p, pow+1)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
