package m68000

import "github.com/f1gopher/gbpixellib/hardware/memory"

type rte struct{}

func CreateRTE() *rte {
	return &rte{}
}

func (r *rte) Name() string {
	return "rte"
}

func (r *rte) Length() int {
	return 2
}

func (r *rte) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
