package cps1

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/f1gopher/gbpixellib/hardware/m68000"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/system/cps1/cartridge"
	"github.com/f1gopher/gbpixellib/system/cps1/memoryBus"
	"github.com/f1gopher/gbpixellib/system/gameboy/display"
)

type CPS1 struct {
	log       *log.Log
	romFolder string

	bus   *memoryBus.Bus
	cpu   *m68000.M68000
	jamma *jamma
}

func CreateCPS1(romFolder string, useDebugger bool) *CPS1 {
	l := log.CreateLog("./log.txt")
	bus, jammaMemory, kickMemory := memoryBus.CreateBus(l)
	jamma := createJamma(l, jammaMemory, kickMemory)

	system := CPS1{
		log:       l,
		romFolder: romFolder,
		bus:       bus,
		cpu:       m68000.CreateM68000(l, bus.Memory68000()),
		jamma:     jamma,
	}

	system.Reset()

	return &system
}

func (c *CPS1) LoadGame(romFolder string) {
	c.romFolder = romFolder
	c.Reset()
}

func (c *CPS1) Start() {
	game, err := cartridge.LoadGame(c.log, c.romFolder, c.bus.Memory68000(), c.bus.MemoryGfx())
	if err != nil {

	}

	c.jamma.InitForGame(game)
}

func (c *CPS1) Reset() {
	c.bus.Reset()
	c.cpu.Reset()

	c.Start()
}

func (c *CPS1) Render(callback func(x int, y int, color color.RGBA)) {
}

func (c *CPS1) SingleFrame() (breakpoint bool, cyclesCompleted uint, err error) {
	return false, 0, nil
}

func (c *CPS1) SingleInstruction() (breakpoint bool, cyclesCompleted uint, err error) {
	c.cpu.ExecuteCycle()
	return false, 0, nil
}

func (c *CPS1) Dump() string {
	var sb strings.Builder

	data, _ := c.bus.Memory68000().DumpCode()

	for x := 0; x < len(data); x += 16 {
		sb.WriteString(fmt.Sprintf("%06X:  %02X%02X %02X%02X %02X%02X %02X%02X %02X%02X %02X%02X %02X%02X %02X%02X  %c%c%c%c%c%c%c%c%c%c%c%c%c%c%c%c\n",
			x,
			data[x],
			data[x+1],
			data[x+2],
			data[x+3],
			data[x+4],
			data[x+5],
			data[x+6],
			data[x+7],
			data[x+8],
			data[x+9],
			data[x+10],
			data[x+11],
			data[x+12],
			data[x+13],
			data[x+14],
			data[x+15],
			displayChar(data[x]),
			displayChar(data[x+1]),
			displayChar(data[x+2]),
			displayChar(data[x+3]),
			displayChar(data[x+4]),
			displayChar(data[x+5]),
			displayChar(data[x+6]),
			displayChar(data[x+7]),
			displayChar(data[x+8]),
			displayChar(data[x+9]),
			displayChar(data[x+10]),
			displayChar(data[x+11]),
			displayChar(data[x+12]),
			displayChar(data[x+13]),
			displayChar(data[x+14]),
			displayChar(data[x+15])))
	}

	return sb.String()
}

func (c *CPS1) LoadTestROM(rom string) {
}

// func (c *CPS1) Joypad() Joypad {
// 	return nil
// }

func (c *CPS1) DisplayConfig() display.DisplayConfig {
	return display.DisplayConfig{
		Width:  320,
		Height: 224,
	}
}

func (c *CPS1) BreakpointReason() string {
	return ""
}

func (c *CPS1) Name() string {
	return "CPS-1"
}

func (c *CPS1) DumpCallstack() []string {
	return make([]string, 0)
}

func displayChar(value uint8) rune {
	if value < 32 || value > 126 {
		return rune(46)
	}
	return rune(int(value))
}
