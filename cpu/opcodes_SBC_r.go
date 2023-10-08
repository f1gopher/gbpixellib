package cpu

import (
	"errors"
	"fmt"
)

type opcode_SBC_r struct {
	opcodeBase

	src register
}

func createSBC_r(opcode uint8, reg register) *opcode_SBC_r {
	return &opcode_SBC_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SBC %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_SBC_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		result, bit3Carry, bit7Carry := subtract8BitAndCarryWithCarry(reg.Get8(A), reg.Get8(o.src), reg.GetFlag(CFlag))
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, bit3Carry)
		reg.setFlag(CFlag, bit7Carry)

		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
