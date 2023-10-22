package cpu

import (
	"errors"
	"fmt"
)

type opcode_CP_r struct {
	opcodeBase

	src Register
}

func createCP_r(opcode uint8, reg Register) *opcode_CP_r {
	return &opcode_CP_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("CP %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_CP_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		result, bit3Carry, bit7Carry := subtract8BitWithCarry(reg.Get8(A), reg.Get8(o.src))
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, true)
		reg.SetFlag(HFlag, bit3Carry)
		reg.SetFlag(CFlag, bit7Carry)

		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
