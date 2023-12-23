package cpu

import (
	"errors"
	"fmt"
)

type opcode_CB_RES_b_r struct {
	opcodeBase

	bit    uint8
	target Register
}

func createCB_RES_b_r(opcode uint8, bit uint8, reg Register) *opcode_CB_RES_b_r {
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

func (o *opcode_CB_RES_b_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.SetRegBit(o.target, o.bit, false)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
