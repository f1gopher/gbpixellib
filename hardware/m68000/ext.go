package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type ext struct {
	opmode   byte
	size     instSize
	register byte
}

func CreateEXT(value uint16) *ext {
	reg := valueFrom3Bits(value, 0)
	opmode := valueFrom3Bits(value, 6)
	size := wordSize
	if opmode == 0b00000011 {
		size = longSize
	}
	return &ext{
		opmode:   opmode,
		size:     size,
		register: reg,
	}
}

func (e *ext) Name() string {
	return fmt.Sprintf("ext.%s   D%d", e.size.string(), e.register)
}

func (e *ext) Length() int {
	return 2
}

func (e *ext) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
