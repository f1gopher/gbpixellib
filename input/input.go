package input

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

type Input struct {
	memory inputMemory
}

func CreateInput(memory inputMemory) *Input {
	return &Input{
		memory: memory,
	}
}

func (i *Input) Reset() {
	// i.memory.WriteByte(0xFF00, 0x3F)
}

func (i *Input) InputStart(value bool) {
	i.setInput(P13, P15, value)
}

func (i *Input) InputSelect(value bool) {
	i.setInput(P12, P15, value)
}

func (i *Input) InputA(value bool) {
	i.setInput(P10, P15, value)
}

func (i *Input) InputB(value bool) {
	i.setInput(P11, P15, value)
}

func (i *Input) InputUp(value bool) {
	i.setInput(P12, P14, value)
}

func (i *Input) InputDown(value bool) {
	i.setInput(P13, P14, value)
}

func (i *Input) InputLeft(value bool) {
	i.setInput(P11, P14, value)
}

func (i *Input) InputRight(value bool) {
	i.setInput(P10, P14, value)
}

func (i *Input) setInput(port1 uint8, port2 uint8, value bool) {
	i.memory.WriteBit(0xFF00, port1, value)
	i.memory.WriteBit(0xFF00, port2, value)
}
