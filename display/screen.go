package display

import "go-boy/memory"

const screenWidth = 160
const screenHeight = 144

const lcdcRegister = 0xFF40
const lcdStatus = 0xFF41

type lcdStatusMode int

const (
	hblank lcdStatusMode = iota
	vblank
	searchOAM
	transferringToController
)

type screenColor int

const (
	white screenColor = iota
	lightGray
	darkGray
	black
)

func (s screenColor) String() string {
	return [...]string{"#", "o", ".", " "}[s]
}

type Screen struct {
	memory *memory.Memory

	buffer []screenColor
}

func CreateScreen(memory *memory.Memory) *Screen {
	return &Screen{
		memory: memory,
		buffer: make([]screenColor, screenWidth*screenHeight),
	}
}

func (s *Screen) Render() string {
	var output string
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			output += s.buffer[y*screenWidth+x].String()
		}
		output += "\n"
	}

	return output
}

func (s *Screen) read() {

}
