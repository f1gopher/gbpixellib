package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type trap struct {
	vector byte
}

func CreateTRAP(value uint16) *trap {
	return &trap{
		vector: byte(value & 0b0000000000001111),
	}
}

func (t *trap) Name() string {
	return fmt.Sprintf("trap    #$%x", t.vector)
}

func (t *trap) Length() int {
	return 2
}

func (t *trap) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
