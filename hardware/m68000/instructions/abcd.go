package instructions

import "fmt"

type abcd struct {
	length     int
	xReg       byte
	yReg       byte
	isMemoryOp bool
}

func CreateABCD(value uint16, code []uint8, codeOffset int) *abcd {
	xReg := valueFrom3Bits(value, 9)
	yReg := valueFrom3Bits(value, 0)
	isMemoryOp := value&0b0000000000001000 == 0b0000000000001000

	return &abcd{
		xReg:       xReg,
		yReg:       yReg,
		isMemoryOp: isMemoryOp,
	}
}

func (a *abcd) Name() string {

	if a.isMemoryOp {
		return fmt.Sprintf("abcd    -(A%d), -(A%d)", a.yReg, a.xReg)
	}

	return fmt.Sprintf("abcd    D%d, D%d", a.yReg, a.xReg)
}

func (a *abcd) Length() int {
	return 2
}
