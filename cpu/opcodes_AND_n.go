package cpu

import (
	"errors"
)

type opcode_AND_n struct {
	opcodeBase

	n uint8
}

func createAND_n(opcode uint8) *opcode_AND_n {
	return &opcode_AND_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "AND n",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_AND_n) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		result := reg.Get8(A) & o.n
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, true)
		reg.setFlag(CFlag, false)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
