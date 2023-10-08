package cpu

import (
	"errors"
)

type opcode_CP_n struct {
	opcodeBase

	n uint8
}

func createCP_n(opcode uint8) *opcode_CP_n {
	return &opcode_CP_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "CP n",
			opcodeLength: 2,
		},
	}
}

func (o *opcode_CP_n) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		result, bit3Carry, bit7Carry := subtract8BitWithCarry(reg.Get8(A), o.n)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, true)
		reg.setFlag(HFlag, bit3Carry)
		reg.setFlag(CFlag, bit7Carry)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
