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

func (i *Input) PressStart() {
	i.inputStart(true)
}

func (i *Input) ReleaseStart() {
	i.inputStart(false)
}

func (i *Input) PressSelect() {
	i.inputSelect(true)
}

func (i *Input) ReleaseSelect() {
	i.inputSelect(false)
}
func (i *Input) PressA() {
	i.inputA(true)
}

func (i *Input) ReleaseA() {
	i.inputA(false)
}
func (i *Input) PressB() {
	i.inputB(true)
}

func (i *Input) ReleaseB() {
	i.inputB(false)
}
func (i *Input) PressUp() {
	i.inputUp(true)
}

func (i *Input) ReleaseUp() {
	i.inputUp(false)
}
func (i *Input) PressDown() {
	i.inputDown(true)
}

func (i *Input) ReleaseDown() {
	i.inputDown(false)
}

func (i *Input) PressLeft() {
	i.inputLeft(true)
}

func (i *Input) ReleaseLeft() {
	i.inputLeft(false)
}
func (i *Input) PressRight() {
	i.inputRight(true)
}
func (i *Input) ReleaseRight() {
	i.inputRight(false)
}
func (i *Input) ReadDirectional() uint8 {
	return i.directional
}

func (i *Input) ReadStandard() uint8 {
	return i.standard
}

// NOTE: 0 is the button is pressed and 1 means not pressed

func (i *Input) inputStart(pressed bool) {
	i.standard = memory.SetBit(i.standard, P13, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputSelect(pressed bool) {
	i.standard = memory.SetBit(i.standard, P12, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputA(pressed bool) {
	i.standard = memory.SetBit(i.standard, P10, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputB(pressed bool) {
	i.standard = memory.SetBit(i.standard, P11, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputUp(pressed bool) {
	i.directional = memory.SetBit(i.directional, P12, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputDown(pressed bool) {
	i.directional = memory.SetBit(i.directional, P13, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputLeft(pressed bool) {
	i.directional = memory.SetBit(i.directional, P11, !pressed)
	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}

func (i *Input) inputRight(pressed bool) {
	i.directional = memory.SetBit(i.directional, P10, !pressed)

	if !pressed {
		i.interupt.Request(interupt.Joypad)
	}
}
