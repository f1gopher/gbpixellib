package cpu

import (
	"errors"
)

type opcode_DEC_abs_HL struct {
	opcodeBase

	data uint8
}

func createDEC_abs_HL(opcode uint8) *opcode_DEC_abs_HL {
	return &opcode_DEC_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "DEC (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_DEC_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.data = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		result, carryBit3, _ := subtract8BitWithCarry(o.data, 1)

		mem.WriteByte(reg.Get16(HL), result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, true)
		reg.SetFlag(HFlag, carryBit3)
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
