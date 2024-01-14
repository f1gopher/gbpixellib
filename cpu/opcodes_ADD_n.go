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
			opcodeLength: 2,
		},
	}
}

func (o *opcode_ADD_n) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		result, bit3Carry, bit7Carry := add8BitWithCarry(reg.Get8(A), o.n)
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, bit3Carry)
		reg.SetFlag(CFlag, bit7Carry)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
