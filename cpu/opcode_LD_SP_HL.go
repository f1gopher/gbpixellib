package cpu

import (
	"errors"
)

type opcode_LD_SP_HL struct {
	opcodeBase

	src  Register
	dest Register
}

func createLD_SP_HL(opcode uint8) *opcode_LD_SP_HL {
	return &opcode_LD_SP_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD SP,HL",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_LD_SP_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.Set16(SP, reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
