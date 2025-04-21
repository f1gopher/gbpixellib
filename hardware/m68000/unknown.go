package m68000

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

func (u *unknown) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
