package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type unknown struct {
	value uint16
}

func CreateUNKNOWN(value uint16) *unknown {
	return &unknown{value: value}
}

func (u *unknown) Name() string {
	return fmt.Sprintf("unknown instruction: %x", u.value)
}

func (u *unknown) Length() int {
	return 2
}

func (u *unknown) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
