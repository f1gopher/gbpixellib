package cpu

import (
	"errors"
	"fmt"
)

type opcode_Invalid struct {
	opcodeBase
}

func createInvalid(opcode uint8) *opcode_Invalid {
	return &opcode_Invalid{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "**Invalid**",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_Invalid) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {
	return true, errors.New(fmt.Sprintf("Invalid opcode 0x%02X", o.opcodeId))
}
