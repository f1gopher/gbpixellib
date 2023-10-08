package cpu

import (
	"errors"
	"fmt"
)

type opcode_LD_A_abs_rr struct {
	opcodeBase

	target register
}

func createLD_A_abs_rr(opcode uint8, reg register) *opcode_LD_A_abs_rr {
	return &opcode_LD_A_abs_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("LD A,(%s)", reg.String()),
			opcodeLength: 1,
		},
		target: reg,
	}
}

func (o *opcode_LD_A_abs_rr) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.set8(A, mem.ReadByte(reg.Get16(o.target)))
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
