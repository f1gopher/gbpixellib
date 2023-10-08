package cpu

import (
	"errors"
)

type opcode_OR_abs_HL struct {
	opcodeBase

	n uint8
}

func createOR_abs_HL(opcode uint8) *opcode_OR_abs_HL {
	return &opcode_OR_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "OR (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_OR_abs_HL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		result := reg.Get8(A) | o.n
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, false)
		reg.setFlag(CFlag, false)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
