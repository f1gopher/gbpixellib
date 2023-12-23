package cpu

import (
	"errors"
)

type opcode_LDH_imed_n_A struct {
	opcodeBase

	n uint8
}

func createLDH_imed_n_A(opcode uint8) *opcode_LDH_imed_n_A {
	return &opcode_LDH_imed_n_A{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LDH (n),A",
			opcodeLength: 2,
		},
	}
}

func (o *opcode_LDH_imed_n_A) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		mem.WriteByte(combineBytes(0xFF, o.n), reg.Get8(A))
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
