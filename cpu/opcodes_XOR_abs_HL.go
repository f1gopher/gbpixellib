package cpu

import (
	"errors"
)

type opcode_XOR_abs_HL struct {
	opcodeBase

	n uint8
}

func createXOR_abs_HL(opcode uint8) *opcode_XOR_abs_HL {
	return &opcode_XOR_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "XOR (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_XOR_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		result := reg.Get8(A) ^ o.n
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, false)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
