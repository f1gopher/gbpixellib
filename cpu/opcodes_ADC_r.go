package cpu

import (
	"errors"
	"fmt"
)

type opcode_ADC_r struct {
	opcodeBase

	src register
}

func createADC_r(opcode uint8, reg register) *opcode_ADC_r {
	return &opcode_ADC_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("ADC %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_ADC_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		result, bit3Carry, bit7Carry := add8BitAndCarryWithCarry(reg.Get8(A), reg.Get8(o.src), reg.GetFlag(CFlag))
		reg.set8(A, result)
		reg.setFlag(ZFlag, result == 0)
		reg.setFlag(NFlag, false)
		reg.setFlag(HFlag, bit3Carry)
		reg.setFlag(CFlag, bit7Carry)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
