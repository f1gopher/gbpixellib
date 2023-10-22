package cpu

import (
	"errors"
	"fmt"
)

type opcode_LD_r_abs_HL struct {
	opcodeBase

	target Register
	value  uint8
}

func createLD_r_abs_HL(opcode uint8, reg Register) *opcode_LD_r_abs_HL {
	return &opcode_LD_r_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("LD %s,(HL)", reg.String()),
			opcodeLength: 1,
		},
		target: reg,
	}
}

func (o *opcode_LD_r_abs_HL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		reg.Set8(o.target, o.value)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
