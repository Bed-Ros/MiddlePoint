package main

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 800

	borderPointsNum = 20
)

var (
	centerColor = color.White
	borderColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	middleColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}

	pointRadius     float32 = 4.0
	borderLineWidth float32 = 1.0
)

type point struct {
	X float32
	Y float32
}

type Game struct {
	points      []point
	middlePoint point
}

func (g *Game) NewDots() {
	g.points = nil

	for i := 0; i < borderPointsNum; i++ {
		newY := rand.Float32() * screenHeight
		newX := rand.Float32() * screenWidth
		g.points = append(g.points, point{X: newX, Y: newY})
	}

	g.UpdateMiddlePoint()
}

func (g *Game) UpdateMiddlePoint() {
	var sumX, sumY float32
	for _, p := range g.points {
		sumY += p.Y
		sumX += p.X
	}

	g.middlePoint = point{
		X: sumX / float32(len(g.points)),
		Y: sumY / float32(len(g.points)),
	}
}

func (g *Game) AddDot(p point) {
	g.points = append(g.points, p)

	g.UpdateMiddlePoint()
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.NewDots()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		g.points = nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
			g.AddDot(point{X: float32(x), Y: float32(y)})
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	//Центр экрана
	//vector.DrawFilledCircle(screen, screenWidth/2, screenHeight/2, pointRadius, centerColor, true)

	//Линии границы
	if len(g.points) >= 2 {
		vector.StrokeLine(screen,
			g.points[0].X, g.points[0].Y,
			g.points[len(g.points)-1].X, g.points[len(g.points)-1].Y,
			borderLineWidth, borderColor, true)
		for i := 0; i < len(g.points)-1; i++ {
			vector.StrokeLine(screen,
				g.points[i].X, g.points[i].Y,
				g.points[i+1].X, g.points[i+1].Y,
				borderLineWidth, borderColor, true)
		}
	}
	//Точки границы
	for _, p := range g.points {
		vector.DrawFilledCircle(screen, p.X, p.Y, pointRadius, borderColor, true)
	}

	//Середина границы
	vector.DrawFilledCircle(screen, g.middlePoint.X, g.middlePoint.Y, pointRadius, middleColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Средняя точка")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
