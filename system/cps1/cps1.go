package cps1

import (
	"image/color"

	"github.com/f1gopher/gbpixellib/hardware/m68000"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/system/cps1/cartridge"
)

type CPS1 struct {
	log       *log.Log
	romFolder string

	bus *bus
	cpu *m68000.M68000
}

func CreateCPS1(romFolder string, useDebugger bool) *CPS1 {
	l := log.CreateLog("./log.txt")
	system := CPS1{
		log:       l,
		romFolder: romFolder,
		bus:       CreateBus(l),
		cpu:       m68000.CreateM68000(l),
	}

	system.Reset()

	return &system
}

func (c *CPS1) LoadGame(romFolder string) {
	c.romFolder = romFolder
	c.Reset()
}

func (c *CPS1) Start() {
	err := cartridge.LoadGame(c.log, c.romFolder, c.bus)
	if err != nil {

	}
}

func (c *CPS1) Reset() {
	c.bus.Reset()

	c.Start()
}

func (c *CPS1) Render(callback func(x int, y int, color color.RGBA)) {
}

func (c *CPS1) SingleFrame() (breakpoint bool, err error) {
	return false, nil
}

func (c *CPS1) SingleInstruction() (breakpoint bool, err error) {
	return false, nil
}
