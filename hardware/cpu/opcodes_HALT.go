package cpu

import (
	"errors"
)

type opcode_HALT struct {
	opcodeBase
}

func createHALT(opcode uint8) *opcode_HALT {
	return &opcode_HALT{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "HALT",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_HALT) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.SetHALT(true)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
