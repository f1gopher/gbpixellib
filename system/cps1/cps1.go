package cps1

import (
	"image/color"
	"log"

	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/system/cartridge"
)

type CPS1 struct {
	log       *log.Log
	romFolder string
}

func CreateCPS1(romFolder string, useDebugger bool) *CPS1 {
	l := log.CreateLog("./log.txt")
	system := CPS1{
		log:       l,
		romFolder: romFolder,
	}

	system.Reset()

	return &system
}

func (c *CPS1) LoadGame(romFolder string) {
	c.romFolder = romFolder
	c.Reset()
}

func (c *CPS1) Start() {
}

func (c *CPS1) Reset() {
	err := cartridge.LoadGame(c.log, c.romFolder)
	if err != nil {

	}
}

func (c *CPS1) Render(callback func(x int, y int, color color.RGBA)) {
}

func (c *CPS1) SingleFrame() (breakpoint bool, err error) {
}

func (c *CPS1) SingleInstruction() (breakpoint bool, err error) {
}
