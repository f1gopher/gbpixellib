package display

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/f1gopher/gbpixellib/interupt"
	"github.com/f1gopher/gbpixellib/memory"
	"golang.org/x/image/colornames"
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

type ScreenColor int

const (
	Off ScreenColor = iota
	White
	LightGray
	DarkGray
	Black
)

func (s ScreenColor) String() string {
	return [...]string{" ", " ", ".", "o", "#"}[s]
}

type interuptHandler interface {
	Request(i interupt.Interupt)
}

type Screen struct {
	log             *os.File
	memory          *memory.Bus
	interuptHandler interuptHandler

	buffer []ScreenColor

	currentCycleForScanline uint
}

func CreateScreen(memory *memory.Bus, interuptHandler interuptHandler) *Screen {
	f, _ := os.Create("./gpu-log.txt")
	return &Screen{
		log:                     f,
		memory:                  memory,
		interuptHandler:         interuptHandler,
		buffer:                  make([]ScreenColor, screenWidth*screenHeight),
		currentCycleForScanline: 0,
	}
}

func (s *Screen) Reset() {
	for x := 0; x < len(s.buffer); x++ {
		s.buffer[x] = Off
	}
	s.currentCycleForScanline = 0
}

// TODO - need to update/set LYC LY compare and LCD status interrupts?

func (s *Screen) DumpTileset() image.Image {
	// A tileset contains 255 tiles each 8x8 pixels
	//
	// Split into two blocks one for 0-127 and one for 128-255 tiles
	// Four rows for each tiles set and a blank line surrounding each tile
	// 10 rows padding between blocks
	// width = 32 * (8 + 1) + 1 =
	// height = 8 * (8 + 1) + 1 + 10 =
	//
	// 0 --------
	// 1
	// 2
	// 3
	// 4
	// 5  1
	// 6
	// 7
	// 8
	// 9 -------
	// 10
	// 17  2
	// 18 -----
	// 19
	// 26  3
	// 27 ------
	// 28
	// 35   4
	// 36 ------
	// 37 ------
	// 45 ------
	// 46  1
	// 53
	// 54 -----
	// 55   2
	// 62 -----
	// 63  3
	// 70
	// 71 -----
	// 72
	// 79   4
	// 80 -----

	img := image.NewRGBA(image.Rect(
		0,
		0,
		32*(8+1)+1,
		(8*8)+18))

	// Border lines are red
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Red}, image.Point{X: 0, Y: 0}, draw.Src)

	tileData := s.BgWindowTileDataArea()
	tileX := 0
	tileY := 0
	var tileNum uint16 = 0
	var tileRow uint16 = 0
	if tileData == 0x8800 {
		tileNum = 255
	}

	for y := 1; y < img.Bounds().Max.Y-1; y++ {

		// Gap between tilesets
		if y > 36 && y <= 45 {
			continue
		}

		// Skip tile seperation rows
		if y == 9 ||
			y == 18 ||
			y == 27 ||
			y == 36 ||
			y == 54 ||
			y == 62 ||
			y == 71 {

			tileY = 0
			tileRow++
			continue
		}

		tileNum = 32 * tileRow
		tileX = 0

		for x := 1; x < img.Bounds().Max.X-1; x++ {

			// Leave a border between tiles
			if tileX == 8 {
				tileX = 0
				if tileData == 0x8000 {
					tileNum++
				} else {
					tileNum--
				}
				continue
			}

			tileAddres := s.tileNumberToAddress(tileData, tileNum, tileY)

			tile := s.memory.ReadShort(tileAddres)

			var colourBit int = int(tileX % 8)
			colourBit -= 7
			colourBit = colourBit * -1

			co := s.colorForPixel(tile, byte(colourBit))

			//// if off then the pixel is transparent
			//if co == White {
			//	continue
			//}

			var c color.RGBA
			switch co {
			case Off:
				c = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			case White:
				c = color.RGBA{R: 155, G: 188, B: 15, A: 255}
			case LightGray:
				c = color.RGBA{R: 139, G: 172, B: 15, A: 255}
			case DarkGray:
				c = color.RGBA{R: 48, G: 98, B: 48, A: 255}
			case Black:
				c = color.RGBA{R: 15, G: 56, B: 15, A: 255}
			default:
				panic("")
			}

			img.Set(x, y, c)

			tileX++
		}

		tileY++
	}

	return img
}

func (s *Screen) tileNumberToAddress(tileData uint16, tileNum uint16, tileY int) uint16 {
	tileLocation := tileData

	if tileData == 0x8000 {
		tileLocation += tileNum * tileSize
	} else {
		if tileNum <= 127 {
			tileLocation = 0x9000
			tileLocation += tileNum * tileSize
		} else {
			tileLocation += (tileNum - 128) * tileSize
		}
	}

	line := tileY % 8
	line = line * 2

	tileAbc := uint16(tileLocation) + uint16(line)
	return tileAbc
}

func (s *Screen) DumpFirstTileMap() *[1024]byte {
	return s.dumpTileMap(0x9800)
}

func (s *Screen) DumpSecondTileMap() *[1024]byte {
	return s.dumpTileMap(0x9C00)
}

func (s *Screen) DumpWindowTileMap() *[1024]byte {
	return s.dumpTileMap(s.WindowTileMapStart())
}

func (s *Screen) DumpBackgroundTileMap() *[1024]byte {
	return s.dumpTileMap(s.BackgroundTileMapStart())
}

func (s *Screen) dumpTileMap(address uint16) *[1024]byte {
	tileMap := [1024]byte{}

	for x := 0; x < 1024; x++ {
		tileNum := s.memory.ReadByte(address + uint16(x))

		tileMap[x] = tileNum
	}

	return &tileMap
}

func (s *Screen) Debug() string {
	return fmt.Sprintf("LCD On: %t\nScrollX: %d\nScrollY: %d\nWindow Enable: %t\nOBJ/Sprite Enabled: %t\nBG Display: %t\nScanline: %d\n",
		s.LCDEnable(),
		s.SCX(),
		s.SCY(),
		s.WindowEnable(),
		s.ObjEnable(),
		s.BgWindowEnablePriority(),
		s.LY(),
	)
}

func (s *Screen) UpdateForCycles(cyclesCompleted uint) {

	s.setLcdMode()

	s.currentCycleForScanline += cyclesCompleted

	if !s.LCDEnable() {
		return
	}

	if s.currentCycleForScanline >= cyclesToDrawScanline {
		currentScanline := s.LY()

		if currentScanline == 144 {
			s.interuptHandler.Request(interupt.VBlank)
		} else if currentScanline > 153 {
			s.memory.DisplaySetScanline(0)
		} else if currentScanline < 144 {
			s.drawScanline()
		}

		// Set bit 2 of stat register to 1
		s.setLcdStatus()

		s.currentCycleForScanline -= cyclesToDrawScanline

		currentScanline = s.LY() + 1
		s.memory.DisplaySetScanline(currentScanline)
	}
}

func (s *Screen) Cycles() uint {
	return s.currentCycleForScanline
}

func (s *Screen) setLcdMode() {
	status := s.memory.ReadByte(lcdStatus)

	currentLine := s.memory.ReadByte(lcdScanline)

	if currentLine > 144 {
		// VBlank - 1
		status = memory.SetBit(status, 0, true)
		status = memory.SetBit(status, 1, false)

	} else {
		if s.currentCycleForScanline <= 80 {
			// Searching OAM Mode - 2
			status = memory.SetBit(status, 0, false)
			status = memory.SetBit(status, 1, true)
		} else if s.currentCycleForScanline >= 167 && s.currentCycleForScanline <= 291 {
			// Transfering Data - 3
			status = memory.SetBit(status, 0, true)
			status = memory.SetBit(status, 1, true)
		} else {
			// HBlank Mode - 0
			status = memory.SetBit(status, 0, false)
			status = memory.SetBit(status, 1, false)

			// Bit 3 is hblank interrupt request
			reqInt := memory.GetBit(status, 3)
			if reqInt {
				s.interuptHandler.Request(interupt.LCD)
			}
		}
	}

	s.memory.DisplaySetStatus(status)
}

func (s *Screen) setLcdStatus() {
	status := s.memory.ReadByte(lcdStatus)

	// LCD interrupt
	a := s.LY()
	b := s.memory.ReadByte(0xFF45)
	bob := a == b

	if bob { //s.LY() == s.memory.ReadByte(0xFF45) {
		status = memory.SetBit(status, 2, true)

		if memory.GetBit(status, 6) {
			s.interuptHandler.Request(interupt.LCD)
		}
	} else {
		status = memory.SetBit(status, 2, false)
	}

	s.memory.DisplaySetStatus(status)
}

func (s *Screen) drawScanline() {
	if s.BgWindowEnablePriority() {
		s.renderTiles()
	}

	if s.ObjEnable() {
		s.renderSprites()
	}
}

func (s *Screen) renderTiles() {

	var backgroundMemory uint16 = 0

	scrollY := s.SCY()
	scrollX := s.SCX()
	windowY := s.WY()
	windowX := s.WX() - 7

	usingWindow := false

	if s.WindowEnable() {
		usingWindow = windowY <= s.LY()
	}

	tileData := s.BgWindowTileDataArea()

	if !usingWindow {
		backgroundMemory = s.BackgroundTileMapStart()
	} else {
		backgroundMemory = s.WindowTileMapStart()
	}

	var yPos byte = 0

	if !usingWindow {
		yPos = scrollY + s.LY()
	} else {
		yPos = s.LY() - windowY
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

		abc := s.memory.ReadByte(tileAddress)
		tileNum = uint16(abc)

		var colourBit int = int(xPos % 8)
		colourBit -= 7
		colourBit = colourBit * -1

		tileAddres := s.tileNumberToAddress(tileData, tileNum, int(yPos))

		tile := s.memory.ReadShort(tileAddres)

		color := s.colorForPixel(tile, byte(colourBit))

		finalY := s.LY()
		if finalY < 0 || finalY >= screenHeight || pixel < 0 || pixel >= screenWidth {
			panic(fmt.Sprintf("Invalid pixel location %d,%d", pixel, finalY))
		}

		offset := (uint16(finalY) * uint16(screenWidth)) + uint16(pixel)

		s.buffer[offset] = color
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

		ySize := s.ObjSize()
		scanline := s.LY()

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

			tile := s.memory.ReadShort(dataAddress)

			for tilePixel := 7; tilePixel >= 0; tilePixel-- {

				colorBit := tilePixel
				if xFlip {
					colorBit -= 7
					colorBit *= -1
				}

				color := s.colorForPixel(tile, byte(colorBit))

				// White is transparent for sprites
				if color == White {
					continue
				}

				xPix := 0 - tilePixel
				xPix += 7
				pixel := uint16(xPos) + uint16(xPix)

				//colourNum := memory.GetBit(data2, int(colourBit))
				//colourNum = colourNum || memory.GetBit(data1, int(colourBit))

				//color := s.bgpColor(colourNum)

				// This happens sometimes (Mario)
				if scanline < 0 || scanline >= screenHeight || pixel < 0 || pixel >= screenWidth {
					// panic("Invalid pixel location")
					continue
				}

				s.buffer[pixel+(uint16(scanline)*uint16(screenWidth))] = color
			}
		}

	}
}

func (s *Screen) Render(callback func(x int, y int, color ScreenColor)) {

	//s.read()

	//var output strings.Builder
	//output.Grow(101 * screenWidth)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			//output.WriteString(s.buffer[(y*screenWidth)+x].String())
			callback(x, y, s.buffer[(y*screenWidth)+x])
		}
		//output.WriteString("\n")
	}

	//return output.String()
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

func (s *Screen) colorForPixel(block uint16, index byte) ScreenColor {
	// var paletteId byte = 0
	// l1 := s.memory.ReadBit(block, 7-index)
	// l2 := s.memory.ReadBit(block+1, 7-index)
	// if l1 {
	// 	paletteId = 1
	// }
	//
	// if l2 {
	// 	paletteId += 2
	// }
	//
	// c := (byte(s.memory.ReadByte(0xFF47)) >> (paletteId * 2)) & 0x03
	// switch c {
	// case 0:
	// 	return white
	// case 1:
	// 	return lightGray
	// case 2:
	// 	return darkGray
	// case 3:
	// 	return black
	// default:
	// 	panic("")
	// }
	//

	// bit 7 is for pixel 0
	//switch index {
	//case 0:
	//	index = 7
	//case 1:
	//	index = 6
	//case 2:
	//	index = 5
	//case 3:
	//	index = 4
	//case 4:
	//	index = 3
	//case 5:
	//	index = 2
	//case 6:
	//	index = 1
	//case 7:
	//	index = 0
	//default:
	//	panic("")
	//}

	highFlag := block >> (8 + index) & 0x0001
	lowFlag := block >> index & 0x0001

	// TODO - check order of bytes

	if highFlag == 0x00 && lowFlag == 0x00 {
		return White
	} else if highFlag == 0x00 && lowFlag == 0x01 {
		return LightGray
	} else if highFlag == 0x01 && lowFlag == 0x00 {
		return DarkGray
	} else { // aFlag == 0x00 && bFlag == 0x00
		return Black
	}
}
