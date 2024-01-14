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
			opcodeLength: 2,
		},
	}
}

func (o *opcode_AND_n) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		result := reg.Get8(A) & o.n
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, true)
		reg.SetFlag(CFlag, false)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
