package cpu

import (
	"errors"
)

type opcode_LDH_A_imed_n struct {
	opcodeBase

	n uint8
}

func createLDH_A_imed_n(opcode uint8) *opcode_LDH_A_imed_n {
	return &opcode_LDH_A_imed_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LDH A,(n)",
			opcodeLength: 2,
		},
	}
}

func (o *opcode_LDH_A_imed_n) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		reg.Set8(A, mem.ReadByte(combineBytes(0xFF, o.n)))
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
