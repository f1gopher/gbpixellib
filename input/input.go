package input

import (
	"github.com/f1gopher/gbpixellib/interupt"
	"github.com/f1gopher/gbpixellib/memory"
)

const P13 = 3
const P12 = 2
const P11 = 1
const P10 = 0

type inputMemory interface {
	WriteBit(address uint16, bit uint8, value bool)
	WriteByte(address uint16, value uint8)
	ReadBit(address uint16, bit uint8) bool
}

type inputInterupt interface {
	Request(i interupt.Interupt)
}

type Input struct {
	memory   inputMemory
	interupt inputInterupt

	directional uint8
	standard    uint8
}

func CreateInput(memory inputMemory, interrupt inputInterupt) *Input {
	return &Input{
		memory:   memory,
		interupt: interrupt,
	}
}

func (i *Input) Reset() {
	i.directional = 0x0F
	i.standard = 0x0F

	// Uncomment for comparison runs
	//i.directional = 0x00
	//i.standard = 0x00
}

func (i *Input) ReadDirectional() uint8 {
	return i.directional
}

func (i *Input) ReadStandard() uint8 {
	return i.standard
}

// NOTE: 0 is the button is pressed and 1 means not pressed

func (i *Input) InputStart(pressed bool) {
	i.standard = memory.SetBit(i.standard, P13, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputSelect(pressed bool) {
	i.standard = memory.SetBit(i.standard, P12, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputA(pressed bool) {
	i.standard = memory.SetBit(i.standard, P10, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputB(pressed bool) {
	i.standard = memory.SetBit(i.standard, P11, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputUp(pressed bool) {
	i.directional = memory.SetBit(i.directional, P12, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputDown(pressed bool) {
	i.directional = memory.SetBit(i.directional, P13, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputLeft(pressed bool) {
	i.directional = memory.SetBit(i.directional, P11, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) InputRight(pressed bool) {
	i.directional = memory.SetBit(i.directional, P10, !pressed)

	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}
