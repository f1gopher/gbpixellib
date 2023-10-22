package cpu

import (
	"errors"
)

type opcode_SBC_n struct {
	opcodeBase

	n uint8
}

func createSBC_n(opcode uint8) *opcode_SBC_n {
	return &opcode_SBC_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "SBC n",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_SBC_n) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 1 {
		result, bit3Carry, bit7Carry := subtract8BitAndCarryWithCarry(reg.Get8(A), o.n, reg.GetFlag(CFlag))
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, true)
		reg.SetFlag(HFlag, bit3Carry)
		reg.SetFlag(CFlag, bit7Carry)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
