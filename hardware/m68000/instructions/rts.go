package instructions

import "github.com/f1gopher/gbpixellib/hardware/memory"

type rts struct{}

func CreateRTS(value uint16) *rts {
	return &rts{}
}

func (r *rts) Name() string {
	return "rts"
}

func (r *rts) Length() int {
	return 2
}

func (r *rts) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
