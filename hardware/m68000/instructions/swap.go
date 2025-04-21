package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type swap struct {
	reg byte
}

func CreateSWAP(value uint16, code []uint8, codeOffset int) *swap {
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

func (s *swap) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
