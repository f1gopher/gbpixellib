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

func (e *ext) Step(registers *Registers, memory memory.Memory[uint32]) bool {

	result := registers.GetDReg32(e.register)

	// Check the most siginificant bit is set depending on size and then extend if set
	if e.size == wordSize {
		if result&0b00000000000000000000000010000000 == 0b00000000000000000000000010000000 {
			result = result | 0b00000000000000001111111100000000
		} else {
			result = result & 0b11111111111111110000000011111111
		}
	} else if e.size == longSize {
		if result&0b00000000000000001000000000000000 == 0b00000000000000001000000000000000 {
			result = result | 0b11111111111111110000000000000000
		} else {
			result = result & 0b00000000000000001111111111111111
		}
	}

	registers.SetDReg32(e.register, result)
	nzvc_NegZeroClrClr(registers, e.size, result)

	return true
}
