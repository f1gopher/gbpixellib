package cpu

import (
	"fmt"
)

type opcode_CB_RES_b_r struct {
	opcodeBase

	bit    uint8
	target register
}

func createCB_RES_b_r(opcode uint8, bit uint8, reg register) *opcode_CB_RES_b_r {
	return &opcode_CB_RES_b_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RES %d,%s", bit, reg.String()),
			opcodeLength: 2, // TODO - Check this
		},
		bit:    bit,
		target: reg,
	}
}

func (o *opcode_CB_RES_b_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		panic("Invalid cycle")
	}

	reg.setRegBit(o.target, o.bit, false)

	return true, nil
}
