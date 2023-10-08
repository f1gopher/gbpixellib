package cpu

import (
	"errors"
)

type opcode_XOR_n struct {
	opcodeBase

	n uint8
}

func createXOR_n(opcode uint8) *opcode_XOR_n {
	return &opcode_XOR_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "XOR n",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_XOR_n) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 1 {
		result := reg.Get8(A) ^ o.n
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, false)
		reg.setFlag(CFlag, false)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
