package cpu

import (
	"errors"
)

type opcode_LDH_imed_C_A struct {
	opcodeBase

	n uint8
}

func createLDH_imed_C_A(opcode uint8) *opcode_LDH_imed_C_A {
	return &opcode_LDH_imed_C_A{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LDH (C),A",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_LDH_imed_C_A) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		mem.WriteByte(CombineBytes(0xFF, reg.Get8(C)), reg.Get8(A))
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
