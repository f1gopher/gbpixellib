package cpu

import (
	"errors"
)

type opcode_NOP struct {
	opcodeBase
}

func createNOP(opcode uint8) *opcode_NOP {
	return &opcode_NOP{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "NOP",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_NOP) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
