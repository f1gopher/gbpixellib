package display

import "fmt"

const screenWidth = 160
const screenHeight = 144

type ScreenColor int

const (
	Color0 ScreenColor = iota
	Color1
	Color2
	Color3
	Blank
)

func (s ScreenColor) String() string {
	return [...]string{"#", "x", "o", ".", " "}[s]
}

type screen struct {
	buffer []ScreenColor
}

func CreateScreen() *screen {
	return &screen{
		buffer: make([]ScreenColor, screenWidth*screenHeight),
	}
}

func (s *screen) Render() {
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			fmt.Print(s.buffer[x*screenWidth+y].String())
		}
		fmt.Print("\n")
	}
}

func (s *screen) Set(x int, y int, value ScreenColor) {
	s.buffer[screenWidth*x+y] = value
}
