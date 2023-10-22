package cpu

import (
	"errors"
	"fmt"
)

type opcode_INC_rr struct {
	opcodeBase

	target Register
}

func createINC_rr(opcode uint8, reg Register) *opcode_INC_rr {
	return &opcode_INC_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("INC %s", reg.String()),
			opcodeLength: 1, // TODO - check this
		},
		target: reg,
	}
}

func (o *opcode_INC_rr) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		original := reg.Get16(o.target)
		result := add16Bit(original, 1)
		reg.Set16(o.target, result)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
