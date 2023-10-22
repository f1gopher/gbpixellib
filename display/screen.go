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
	memory          *memory.Memory
	interuptHandler interuptHandler

	buffer []ScreenColor

	currentCycleForScanline int
}

func CreateScreen(memory *memory.Memory, interuptHandler interuptHandler) *Screen {
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

	//currentAddr := s.BgWindowTileDataArea()
	tileData := s.BgWindowTileDataArea()
	unsig := tileData == 0x8000
	tileX := 0
	tileY := 0
	var tileNum uint16 = 0
	var tileRow uint16 = 0

	for y := 1; y < img.Bounds().Max.Y-1; y++ {

		// Skip tile seperation rows
		if y == 9 ||
			y == 18 ||
			y == 27 ||
			(y >= 36 && y <= 45) ||
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
				tileNum++
				continue
			}

			tileLocation := tileData

			if unsig {
				tileLocation += tileNum * 0x10
			} else {
				tileLocation += (tileNum + 128) * 0x10
			}

			line := tileY % 8
			line = line * 2
			//data1 := s.memory.ReadByte(tileLocation + line)
			//data2 := s.memory.ReadByte(tileLocation + line + 1)

			var colourBit int = int(tileX % 8)
			colourBit -= 7
			colourBit = colourBit * -1

			tileAbc := uint16(tileLocation) + uint16(line)

			//tileAbc = s.tileAddrForId(abc)

			tile := s.memory.ReadShort(tileAbc)

			co := s.colorForPixel(tile, byte(colourBit))

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

func (s *Screen) DumpTileMap() *[1024]byte {
	tileMap := [1024]byte{}

	// Background tilemap
	backgroundMemory := s.BgTileMapArea(6)

	for x := 0; x < 1024; x++ {
		tileNum := s.memory.ReadByte(backgroundMemory + uint16(x))

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

func (s *Screen) UpdateForCycles(cyclesCompleted int) {

	s.setLcdStatus()

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

		s.currentCycleForScanline -= cyclesToDrawScanline

		currentScanline = s.LY() + 1
		s.memory.DisplaySetScanline(currentScanline)
	}
}

func (s *Screen) Cycles() int {
	return s.currentCycleForScanline
}

func (s *Screen) setLcdStatus() {
	status := s.memory.ReadByte(lcdStatus)

	//if !s.LCDEnable() {
	//	s.currentCycleForScanline = 0
	//	s.memory.DisplaySetScanline(0)
	//	status = status & 252
	//	status = memory.SetBit(status, 0, true)
	//	s.memory.WriteByte(lcdStatus, status)
	//	return
	//}

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

	if s.LY() == s.memory.ReadByte(0xFF45) {
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

	s.memory.DumpTiles()

	if s.WindowEnable() {
		usingWindow = windowY <= s.LY()
	}

	tileData := s.BgWindowTileDataArea()
	unsig := tileData == 0x8000

	if !usingWindow {
		backgroundMemory = s.BgTileMapArea(3)
	} else {
		backgroundMemory = s.BgTileMapArea(6)
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

		// TODO - tileAddress is wrong??

		// if unsig {
		abc := s.memory.ReadByte(tileAddress)
		tileNum = uint16(abc)
		// } else {
		// 	tileNum = uint16(s.memory.ReadByte(tileAddress))
		// }

		// TODO - HACK
		// tileNum = 1
		// TODO - temp hack to draw something
		//tileNum = 5

		tileLocation := tileData

		if unsig {
			tileLocation += tileNum * 0x10
		} else {
			tileLocation += (tileNum + 128) * 0x10
		}

		line := yPos % 8
		line = line * 2
		//data1 := s.memory.ReadByte(tileLocation + line)
		//data2 := s.memory.ReadByte(tileLocation + line + 1)

		var colourBit int = int(xPos % 8)
		colourBit -= 7
		colourBit = colourBit * -1

		tileAbc := uint16(tileLocation) + uint16(line)

		//tileAbc = s.tileAddrForId(abc)

		tile := s.memory.ReadShort(tileAbc)

		color := s.colorForPixel(tile, byte(colourBit))

		//colourNum := memory.GetBit(data2, int(colourBit))
		//colourNum = colourNum || memory.GetBit(data1, int(colourBit))

		//color := s.bgpColor(colourNum)

		finalY := s.LY()
		if finalY < 0 || finalY >= screenHeight || pixel < 0 || pixel >= screenWidth {
			panic(fmt.Sprintf("Invalid pixel location %d,%d", pixel, finalY))
		}

		offset := (uint16(finalY) * uint16(screenWidth)) + uint16(pixel)

		s.buffer[offset] = color

		if pixel == 32 {
			s.log.WriteString(fmt.Sprintf("Y: %d, tileid: %d, tile addr: 0x%04X, data: 0x%04X\n", tileRow, tileNum, tileAbc, tile))
		}
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

			for tilePixel := 7; tilePixel >= 0; tilePixel-- {

				colorBit := tilePixel
				if xFlip {
					colorBit -= 7
					colorBit *= -1
				}

				color := s.colorForPixel(dataAddress, byte(colorBit))

				// White is transparent for sprites
				if color == White {
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

func (s *Screen) read() {
	currentTileAddr := s.BgTileMapArea(0)

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
