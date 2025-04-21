package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type sbcd struct {
	length     int
	xReg       byte
	yReg       byte
	isMemoryOp bool
}

func CreateSBCD(value uint16, code []uint8, codeOffset int) *sbcd {
	xReg := valueFrom3Bits(value, 9)
	yReg := valueFrom3Bits(value, 0)
	isMemoryOp := value&0b0000000000001000 == 0b0000000000001000

	return &sbcd{
		xReg:       xReg,
		yReg:       yReg,
		isMemoryOp: isMemoryOp,
	}
}

func (s *sbcd) Name() string {

	if s.isMemoryOp {
		return fmt.Sprintf("sbcd    -(A%d), -(A%d)", s.yReg, s.xReg)
	}

	return fmt.Sprintf("sbcd    D%d, D%d", s.yReg, s.xReg)
}

func (s *sbcd) Length() int {
	return 2
}

func (s *sbcd) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
