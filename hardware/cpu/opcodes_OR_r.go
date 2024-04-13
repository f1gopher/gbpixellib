package cpu

import (
	"errors"
	"fmt"
)

type opcode_OR_r struct {
	opcodeBase

	src Register
}

func createOR_r(opcode uint8, reg Register) *opcode_OR_r {
	return &opcode_OR_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("OR %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_OR_r) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		result := reg.Get8(A) | reg.Get8(o.src)
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, false)

		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
