package cpu

import (
	"errors"
)

type opcode_ADC_abs_HL struct {
	opcodeBase

	n uint8
}

func createADC_abs_HL(opcode uint8) *opcode_ADC_abs_HL {
	return &opcode_ADC_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "ADC (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_ADC_abs_HL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 1 {
		result, bit3Carry, bit7Carry := add8BitAndCarryWithCarry(reg.Get8(A), o.n, reg.GetFlag(CFlag))
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, bit3Carry)
		reg.setFlag(CFlag, bit7Carry)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
