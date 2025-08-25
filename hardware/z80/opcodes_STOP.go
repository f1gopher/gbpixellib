package z80

import (
	"errors"
)

type opcode_STOP struct {
	opcodeBase
}

func createSTOP(opcode uint8) *opcode_STOP {
	return &opcode_STOP{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "STOP",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_STOP) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
