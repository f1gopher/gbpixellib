package cpu

import (
	"errors"
	"fmt"
)

type opcode_LD_rr_nn struct {
	opcodeBase

	target Register
	lsb    uint8
}

func createLD_rr_nn(opcode uint8, reg Register) *opcode_LD_rr_nn {
	return &opcode_LD_rr_nn{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("LD %s,nn", reg.String()),
			opcodeLength: 3,
		},
		target: reg,
	}
}

func (o *opcode_LD_rr_nn) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		msb := readAndIncPC(reg, mem)
		reg.Set16(o.target, combineBytes(msb, o.lsb))
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
