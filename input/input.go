package input

import "github.com/f1gopher/gbpixellib/interupt"

const P15 = 5
const P14 = 4
const P13 = 3
const P12 = 2
const P11 = 1
const P10 = 0

type inputMemory interface {
	WriteBit(address uint16, bit uint8, value bool)
	WriteByte(address uint16, value uint8)
}

type inputInterupt interface {
	Request(i interupt.Interupt)
}

type Input struct {
	memory   inputMemory
	interupt inputInterupt
}

func CreateInput(memory inputMemory, interrupt inputInterupt) *Input {
	return &Input{
		memory:   memory,
		interupt: interrupt,
	}
}

func (i *Input) Reset() {
	// i.memory.WriteByte(0xFF00, 0x3F)
}

func (i *Input) InputStart(value bool) {
	i.setInput(P13)
}

func (i *Input) InputSelect(value bool) {
	i.setInput(P12)
}

func (i *Input) InputA(value bool) {
	i.setInput(P10)
}

func (i *Input) InputB(value bool) {
	i.setInput(P11)
}

func (i *Input) InputUp(value bool) {
	i.setInput(P11)
}

func (i *Input) InputDown(value bool) {
	i.setInput(P13)
}

func (i *Input) InputLeft(value bool) {
	i.setInput(P11)
}

func (i *Input) InputRight(value bool) {
	i.setInput(P10)
}

func (i *Input) setInput(port1 uint8) {
	i.memory.WriteBit(0xFF00, port1, false)
	i.interupt.Request(interupt.Joypad)
}
