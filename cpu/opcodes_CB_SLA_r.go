package cpu

import (
	"fmt"
)

type opcode_CB_SLA_r struct {
	opcodeBase

	target register
}

func createCB_SLA_r(opcode uint8, reg register) *opcode_CB_SLA_r {
	return &opcode_CB_SLA_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SLA %s", reg.String()),
			opcodeLength: 2, // TODO - check this
		},
		target: reg,
	}
}

func (o *opcode_CB_SLA_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		panic("Invalid cycle")
	}

	panic("Not implemented")

	return true, nil
}
