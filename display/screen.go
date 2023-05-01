package display

import (
	"fmt"
	"go-boy/memory"
	"strings"
)

const screenWidth = 160
const screenHeight = 144

const lcdcRegister = 0xFF40
const lcdStatus = 0xFF41

const tileSize = 16

type lcdStatusMode int

const (
	hblank lcdStatusMode = iota
	vblank
	searchOAM
	transferringToController
)

type screenColor int

const (
	off screenColor = iota
	white
	lightGray
	darkGray
	black
)

func (s screenColor) String() string {
	return [...]string{"X", " ", ".", "o", "#"}[s]
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

func (s *Screen) Debug() string {
	return fmt.Sprintf("LCD On: %t\nScrollX: %d\nScrollY: %d\nWindow Enable: %t\nOBJ/Sprite Enabled: %t\nBG Display: %t\n",
		s.lcdEnable(),
		s.scx(),
		s.scy(),
		s.windowEnable(),
		s.objEnable(),
		s.bgWindowEnablePriority(),
	)
}

func (s *Screen) Render() string {

	s.read()

	var output strings.Builder
	output.Grow(101 * screenWidth)
	for y := 0; y < 100; y++ {
		for x := 0; x < screenWidth; x++ {
			output.WriteString(s.buffer[y*screenWidth+x].String())
		}
		output.WriteString("\n")
	}

	return output.String()
}

func (s *Screen) read() {
	currentTileAddr := s.bgTileMapArea()

	//	scrollY := s.scy()
	//	scrollX := s.scx()

	for y := 0; y < screenHeight; y += 8 {
		for x := 0; x < screenWidth; x += 8 {
			tileId := s.memory.ReadByte(currentTileAddr)
			tileAddr := s.tileAddrForId(tileId)

			s.drawTile(tileAddr, (screenWidth*y)+x)

			currentTileAddr++
		}
	}

	// FF40 controls layout in memory

	// draw background tiles - only ones visible in region
	// 32x32 set cropped down
	// bg map

	// draw window
	// window map

	// draw sprites
	// sprites have OAM attributes
	// 40 sprites in system
	// 40 OAM entries in OAM RAM

	// draw 60 times a second
}

func (s *Screen) tileAddrForId(id byte) uint16 {
	lcdc4 := s.memory.ReadBit(lcdcRegister, 4)

	if lcdc4 {
		if id <= 127 {
			return 0x8000 + uint16(id)
		} else {
			return 0x8800 + (uint16(id) - 127)
		}
	} else {
		if id <= 127 {
			return 0x9000 + uint16(id)
		} else {
			return 0x8800 + (uint16(id) - 127)
		}
	}
}

func (s *Screen) drawTile(tileAddr uint16, firstPixelIndex int) {

	pixelIndex := firstPixelIndex

	for pixelBlockAddr := tileAddr; pixelBlockAddr < tileAddr+tileSize; pixelBlockAddr += 2 {
		pixelBlock := s.memory.ReadShort(pixelBlockAddr)

		for x := 0; x < 8; x++ {
			color := s.colorForPixel(pixelBlock, byte(x))

			s.buffer[pixelIndex+x] = color
		}

		// Move to next row of pixels
		pixelIndex += screenWidth
	}
}

func (s *Screen) colorForPixel(block uint16, index byte) screenColor {
	highFlag := block >> (8 + index) & 0x0001
	lowFlag := block >> index & 0x0001

	if highFlag == 0x00 && lowFlag == 0x00 {
		return white
	} else if highFlag == 0x00 && lowFlag == 0x01 {
		return lightGray
	} else if highFlag == 0x01 && lowFlag == 0x00 {
		return darkGray
	} else { // aFlag == 0x00 && bFlag == 0x00
		return black
	}
}
