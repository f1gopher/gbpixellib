package m68000

import "github.com/f1gopher/gbpixellib/hardware/memory"

type pmove struct {
	name string
}

func CreatePMOVE(value uint16) *pmove {
	return &pmove{name: "pmove [unknown form] D0"}
}

func (p *pmove) Name() string {
	return p.name
}

func (p *pmove) Length() int {
	return 4
}

func (p *pmove) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
