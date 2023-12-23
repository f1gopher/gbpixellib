package cpu

import (
	"errors"
	"fmt"
)

type opcode_LD_r_r struct {
	opcodeBase

	src  Register
	dest Register
}

func createLD_r_r(opcode uint8, dest Register, src Register) *opcode_LD_r_r {
	return &opcode_LD_r_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("LD %s,%s", dest.String(), src.String()),
			opcodeLength: 1,
		},
		src:  src,
		dest: dest,
	}
}

func (o *opcode_LD_r_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.Set8(o.dest, reg.Get8(o.src))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
