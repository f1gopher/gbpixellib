package cpu

import (
	"errors"
)

type opcode_ADD_n struct {
	opcodeBase

	n uint8
}

func createADD_n(opcode uint8) *opcode_ADD_n {
	return &opcode_ADD_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "ADD n",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_ADD_n) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		result, bit3Carry, bit7Carry := add8BitWithCarry(reg.Get8(A), o.n)
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, bit3Carry)
		reg.setFlag(CFlag, bit7Carry)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
