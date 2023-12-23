package cpu

import (
	"errors"
)

type opcode_SBC_abs_HL struct {
	opcodeBase

	n uint8
}

func createSBC_abs_HL(opcode uint8) *opcode_SBC_abs_HL {
	return &opcode_SBC_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "SBC (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_SBC_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		result, bit4Carry, noBorrow := subtract8BitWithCarryBit4(reg.Get8(A), o.n, reg.GetFlag(CFlag))
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, true)
		reg.SetFlag(HFlag, bit4Carry)
		reg.SetFlag(CFlag, noBorrow)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
