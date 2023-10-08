package cpu

import (
	"errors"
	"fmt"
)

type opcode_LD_r_n struct {
	opcodeBase

	target register
	value  uint8
}

func createLD_r_n(opcode uint8, reg register) *opcode_LD_r_n {
	return &opcode_LD_r_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("LD %s,n", reg.String()),
			opcodeLength: 2,
		},
		target: reg,
	}
}

func (o *opcode_LD_r_n) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		reg.set8(o.target, o.value)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
