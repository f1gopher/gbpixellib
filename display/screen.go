package display

import (
	"fmt"
	"go-boy/interupt"
	"go-boy/memory"
	"strings"
)

const screenWidth = 160
const screenHeight = 144

const lcdcRegister = 0xFF40
const lcdStatus = 0xFF41
const lcdScanline = 0xFF44

const tileSize = 16

const cyclesToDrawScanline = 456

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
	return [...]string{" ", " ", ".", "o", "#"}[s]
}

type interuptHandler interface {
	Request(i interupt.Interupt)
}

type Screen struct {
	memory          *memory.Memory
	interuptHandler interuptHandler

	buffer []screenColor

	currentCycleForScanline int
}

func CreateScreen(memory *memory.Memory, interuptHandler interuptHandler) *Screen {
	return &Screen{
		memory:          memory,
		interuptHandler: interuptHandler,
		buffer:          make([]screenColor, screenWidth*screenHeight),
	}
}

func (s *Screen) Debug() string {
	return fmt.Sprintf("LCD On: %t\nScrollX: %d\nScrollY: %d\nWindow Enable: %t\nOBJ/Sprite Enabled: %t\nBG Display: %t\nScanline: %d\n",
		s.lcdEnable(),
		s.scx(),
		s.scy(),
		s.windowEnable(),
		s.objEnable(),
		s.bgWindowEnablePriority(),
		s.ly(),
	)
}

func (s *Screen) UpdateForCycles(cyclesCompleted int) {

	s.setLcdStatus()

	if !s.lcdEnable() {
		return
	}

	s.currentCycleForScanline -= cyclesCompleted

	if s.currentCycleForScanline <= 0 {
		s.currentCycleForScanline = cyclesToDrawScanline

		currentScanline := s.ly() + 1
		s.memory.DisplaySetScanline(currentScanline)

		if currentScanline == 144 {
			s.interuptHandler.Request(interupt.VBlank)
		} else if currentScanline > 153 {
			s.memory.DisplaySetScanline(0)
		} else if currentScanline < 144 {
			s.drawScanline()
		}
	}
}

func (s *Screen) setLcdStatus() {
	status := s.memory.ReadByte(lcdStatus)

	if !s.lcdEnable() {
		s.currentCycleForScanline = 456
		s.memory.DisplaySetScanline(0)
		status = status & 252
		status = memory.SetBit(status, 0, true)
		s.memory.WriteByte(lcdStatus, status)
		return
	}

	currentLine := s.memory.ReadByte(lcdScanline)
	currentMode := status & 0x3

	var mode byte = 0
	reqInt := false

	if currentLine >= 144 {
		mode = 1
		status = memory.SetBit(status, 0, true)
		status = memory.SetBit(status, 1, false)
		reqInt = memory.GetBit(status, 4)

	} else {
		mode2Bounds := 456 - 80
		mode3Bounds := mode2Bounds - 172

		if s.currentCycleForScanline >= mode2Bounds {
			mode = 2
			status = memory.SetBit(status, 1, true)
			status = memory.SetBit(status, 0, false)
			reqInt = memory.GetBit(status, 5)
		} else if s.currentCycleForScanline >= mode3Bounds {
			mode = 3
			status = memory.SetBit(status, 1, true)
			status = memory.SetBit(status, 0, true)
		} else {
			mode = 0
			status = memory.SetBit(status, 1, false)
			status = memory.SetBit(status, 0, false)
			reqInt = memory.GetBit(status, 3)
		}
	}

	if reqInt && mode != currentMode {
		s.interuptHandler.Request(interupt.LCD)
	}

	if s.ly() == s.memory.ReadByte(0xFF45) {
		status = memory.SetBit(status, 2, true)

		if memory.GetBit(status, 6) {
			s.interuptHandler.Request(interupt.LCD)
		}
	} else {
		status = memory.SetBit(status, 2, false)
	}

	s.memory.WriteByte(lcdStatus, status)
}

func (s *Screen) drawScanline() {
	if s.bgWindowEnablePriority() {
		s.renderTiles()
	}

	if s.objEnable() {
		s.renderSprites()
	}
}

func (s *Screen) renderTiles() {

	var backgroundMemory uint16 = 0

	scrollY := s.scy()
	scrollX := s.scx()
	windowY := s.wy()
	windowX := s.wx() - 7

	usingWindow := false

	s.memory.DumpTiles()

	if s.windowEnable() {
		usingWindow = windowY <= s.ly()
	}

	tileData := s.bgWindowTileDataArea()
	unsig := tileData == 0x8000

	if !usingWindow {
		backgroundMemory = s.bgTileMapArea(3)
	} else {
		backgroundMemory = s.bgTileMapArea(6)
	}

	var yPos byte = 0

	if !usingWindow {
		yPos = scrollY + s.ly()
	} else {
		yPos = s.ly() - windowY
	}

	tileRow := ((uint16(yPos) % 0x100) / 8) * 32

	var pixel byte = 0
	for pixel = 0; pixel < screenWidth; pixel++ {
		xPos := pixel + scrollX

		if usingWindow {
			if pixel >= windowX {
				xPos = pixel - windowX
			}
		}

		tileCol := uint16(xPos) / 8 % 32
		var tileNum uint16 = 0
		tileAddress := backgroundMemory + tileRow + tileCol

		// TODO - tileAddress is wrong??

		if unsig {
			tileNum = uint16(s.memory.ReadByte(tileAddress))
		} else {
			tileNum = uint16(s.memory.ReadByte(tileAddress))
		}

		// TODO - temp hack to draw something
		//tileNum = 5

		tileLocation := tileData

		if unsig {
			tileLocation += tileNum * 16
		} else {
			tileLocation += (tileNum + 128) * 16
		}

		line := yPos % 8
		line = line * 2
		//data1 := s.memory.ReadByte(tileLocation + line)
		//data2 := s.memory.ReadByte(tileLocation + line + 1)

		var colourBit int = int(xPos % 8)
		colourBit -= 7
		colourBit = colourBit * -1

		tile := s.memory.ReadShort(tileLocation + uint16(line))

		color := s.colorForPixel(tile, byte(colourBit))

		//colourNum := memory.GetBit(data2, int(colourBit))
		//colourNum = colourNum || memory.GetBit(data1, int(colourBit))

		//color := s.bgpColor(colourNum)

		finalY := s.ly()
		if finalY < 0 || finalY >= screenHeight || pixel < 0 || pixel >= screenWidth {
			panic("Invalid pixel location")
		}

		s.buffer[pixel+(finalY*screenWidth)] = color
	}
}

func (s *Screen) renderSprites() {

	var sprite uint16 = 0
	for sprite = 0; sprite < 40; sprite++ {

		index := sprite * 4
		yPos := s.memory.ReadByte(0xFE00+index) - 16
		xPos := s.memory.ReadByte(0xFE00+index+1) - 8
		tileLocation := s.memory.ReadByte(0xFE00 + index + 2)
		attributes := s.memory.ReadByte(0xFE00 + index + 3)

		//		yFlip := memory.GetBit(attributes, 6)
		xFlip := memory.GetBit(attributes, 5)

		ySize := s.objSize()
		scanline := s.ly()

		if scanline >= yPos && scanline < (yPos+ySize) {

			line := scanline - yPos

			//if yFlip {
			//	line -= byte(ySize)
			//	line *= -1
			//}

			line *= 2
			dataAddress := (0x8000 + (uint16(tileLocation) * 16)) + uint16(line)
			//data1 := s.memory.ReadByte(dataAddress)
			//data2 := s.memory.ReadByte(dataAddress + 1)

			for tilePixel := 7; tilePixel >= 0; tilePixel-- {

				colorBit := tilePixel
				if xFlip {
					colorBit -= 7
					colorBit *= -1
				}

				color := s.colorForPixel(dataAddress, byte(colorBit))

				// White is transparent for sprites
				if color == white {
					continue
				}

				xPix := 0 - tilePixel
				xPix += 7
				pixel := xPos + byte(xPix)

				//colourNum := memory.GetBit(data2, int(colourBit))
				//colourNum = colourNum || memory.GetBit(data1, int(colourBit))

				//color := s.bgpColor(colourNum)

				if scanline < 0 || scanline >= screenHeight || pixel < 0 || pixel >= screenWidth {
					panic("Invalid pixel location")
				}

				s.buffer[pixel+(scanline*screenWidth)] = color
			}
		}

	}
}

func (s *Screen) Render() string {

	//s.read()

	var output strings.Builder
	output.Grow(101 * screenWidth)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			output.WriteString(s.buffer[y*screenWidth+x].String())
		}
		output.WriteString("\n")
	}

	return output.String()
}

func (s *Screen) read() {
	currentTileAddr := s.bgTileMapArea(0)

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
