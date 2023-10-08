package cpu

import (
	"errors"
)

type opcode_DI struct {
	opcodeBase
}

func createDI(opcode uint8) *opcode_DI {
	return &opcode_DI{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "DI",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_DI) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.SetIME(false)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
