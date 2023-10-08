package cpu

import (
	"errors"
	"fmt"
)

type opcode_OR_r struct {
	opcodeBase

	src register
}

func createOR_r(opcode uint8, reg register) *opcode_OR_r {
	return &opcode_OR_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("OR %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_OR_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		result := reg.Get8(A) | reg.Get8(o.src)
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, false)
		reg.setFlag(CFlag, false)

		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
