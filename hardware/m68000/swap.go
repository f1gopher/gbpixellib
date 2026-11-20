package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type swap struct {
	reg byte
}

func CreateSWAP(value uint16) *swap {
	return &swap{
		reg: valueFrom3Bits(value, 0),
	}
}

func (s *swap) Name() string {
	return fmt.Sprintf("swap    D%d", s.reg)
}

func (s *swap) Length() int {
	return 2
}

func (s *swap) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	result := registers.GetDReg32(s.reg)
	// Swap the upper and lower 16bits
	result = (result << 16) | (result >> 16)
	registers.SetDReg32(s.reg, result)

	nzvc_NegZeroClrClr(registers, longSize, result)
	return true
}
