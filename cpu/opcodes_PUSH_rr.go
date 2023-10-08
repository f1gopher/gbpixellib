package cpu

import (
	"errors"
	"fmt"
)

type opcode_PUSH_rr struct {
	opcodeBase

	target register
	msb    uint8
}

func createPUSH_rr(opcode uint8, reg register) *opcode_PUSH_rr {
	return &opcode_PUSH_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("PUSH %s", reg.String()),
			opcodeLength: 1,
		},
		target: reg,
	}
}

func (o *opcode_PUSH_rr) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {
	if cycleNumber == 1 {
		decSP(reg)
		return false, nil
	}

	if cycleNumber == 2 {
		mem.WriteByte(reg.Get16(SP), reg.get16Msb(o.target))
		decSP(reg)
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(SP), reg.get16Lsb(o.target))
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
