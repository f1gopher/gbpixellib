package cpu

import (
	"errors"
	"fmt"
)

type opcode_POP_rr struct {
	opcodeBase

	target register
	lsb    uint8
}

func createPOP_rr(opcode uint8, reg register) *opcode_POP_rr {
	return &opcode_POP_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("POP %s", reg.String()),
			opcodeLength: 1,
		},
		target: reg,
	}
}

func (o *opcode_POP_rr) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {
	if cycleNumber == 1 {
		o.lsb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		msb := readAndIncSP(reg, mem)
		reg.set16FromTwoBytes(o.target, msb, o.lsb)
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
