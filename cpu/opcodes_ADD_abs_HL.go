package cpu

import (
	"errors"
)

type opcode_ADD_abs_HL struct {
	opcodeBase

	n uint8
}

func createADD_abs_HL(opcode uint8) *opcode_ADD_abs_HL {
	return &opcode_ADD_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "ADD (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_ADD_abs_HL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = mem.ReadByte(reg.Get16(HL))
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
