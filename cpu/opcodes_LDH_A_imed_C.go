package cpu

import (
	"errors"
)

type opcode_LDH_A_imed_C struct {
	opcodeBase

	n uint8
}

func createLDH_A_imed_C(opcode uint8) *opcode_LDH_A_imed_C {
	return &opcode_LDH_A_imed_C{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LDH A,(C)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_LDH_A_imed_C) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.Set8(A, mem.ReadByte(combineBytes(0xFF, reg.Get8(C))))
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
