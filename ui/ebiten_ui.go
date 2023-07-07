package ui

import (
	"fmt"
	"go-boy/display"
	"go-boy/system"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const displayWidth = 160
const displayHeight = 144
const scaling float64 = 2

type EbitenUI struct {
	system *system.System
	buffer *ebiten.Image
}

func (e *EbitenUI) Update() error {
	e.system.Tick()

	e.system.Render(func(x int, y int, pixelColor display.ScreenColor) {
		var c color.RGBA
		switch pixelColor {
		case display.Off:
			c = color.RGBA{255, 0, 0, 255}
		case display.White:
			c = color.RGBA{155, 188, 15, 255}
		case display.LightGray:
			c = color.RGBA{139, 172, 15, 255}
		case display.DarkGray:
			c = color.RGBA{48, 98, 48, 255}
		case display.Black:
			c = color.RGBA{15, 56, 15, 255}
		default:
			panic("")

		}
		e.buffer.Set(x, y, c)
	})

	return nil
}

func (e *EbitenUI) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(e.buffer, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f, TPS: %f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (e *EbitenUI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return displayWidth, displayHeight
}

func (e *EbitenUI) Main(system *system.System) {
	e.system = system
	e.buffer = ebiten.NewImage(displayWidth, displayHeight)

	ebiten.SetWindowSize(int(float64(displayWidth)*scaling), int(float64(displayHeight)*scaling))
	ebiten.SetWindowTitle("GBPixel")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}
