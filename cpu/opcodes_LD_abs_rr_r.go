package cpu

import (
	"errors"
	"fmt"
)

type opcode_LD_abs_rr_r struct {
	opcodeBase

	src  Register
	dest Register
}

func createLD_abs_rr_r(opcode uint8, dest Register, src Register) *opcode_LD_abs_rr_r {
	return &opcode_LD_abs_rr_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("LD (%s),%s", dest.String(), src.String()),
			opcodeLength: 1,
		},
		src:  src,
		dest: dest,
	}
}

func (o *opcode_LD_abs_rr_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		mem.WriteByte(reg.Get16(o.dest), reg.Get8(o.src))
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
