package cpu

import (
	"errors"
)

type opcode_INC_abs_HL struct {
	opcodeBase

	data uint8
}

func createINC_abs_HL(opcode uint8) *opcode_INC_abs_HL {
	return &opcode_INC_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "INC (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_INC_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.data = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		result, carryBit3, _ := add8BitWithCarry(o.data, 1)

		mem.WriteByte(reg.Get16(HL), result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, carryBit3)
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
