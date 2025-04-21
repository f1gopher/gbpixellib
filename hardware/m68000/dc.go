package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type dc struct {
	value uint16
}

func CreateDC(value uint16) *dc {
	return &dc{
		value: value,
	}
}

func (d *dc) Name() string {
	if d.value&0xF000 == 0xF000 {
		return fmt.Sprintf("dc.w    $%x; opcode 1111", d.value)
	}

	return fmt.Sprintf("dc.w    $%04x; ILLEGAL", d.value)
}

func (d *dc) Length() int {
	return 2
}

func (d *dc) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
