package cpu

import (
	"errors"
)

type opcode_EI struct {
	opcodeBase
}

func createEI(opcode uint8) *opcode_EI {
	return &opcode_EI{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "EI",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_EI) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.SetIME(true)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
