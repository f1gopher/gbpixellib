package cpu

import (
	"fmt"
)

type opcode_DEC_rr struct {
	opcodeBase

	target register
}

func createDEC_rr(opcode uint8, reg register) *opcode_DEC_rr {
	return &opcode_DEC_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("DEC %s", reg.String()),
			opcodeLength: 1, // TODO - check this
		},
		target: reg,
	}
}

func (o *opcode_DEC_rr) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		panic("Invalid cycle")
	}

	original := reg.Get16(o.target)
	result := subtract16Bit(original, 1)
	reg.set16(o.target, result)

	return true, nil
}
