package cpu

import (
	"errors"
	"fmt"
)

type opcode_ADD_HL_rr struct {
	opcodeBase

	src Register
}

func createADD_HL_rr(opcode uint8, srcReg Register) *opcode_ADD_HL_rr {
	return &opcode_ADD_HL_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("ADD HL,%s", srcReg.String()),
			opcodeLength: 1, // TODO - check
		},
		src: srcReg,
	}
}

func (o *opcode_ADD_HL_rr) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	// TODO - what happens in which cycle is probably wrong
	if cycleNumber == 1 {
		src := reg.Get16(o.src)
		hl := reg.Get16(HL)

		result, bit11Carry, bit15Carry := add16BitWithCarry(hl, src)

		reg.Set16(HL, result)

		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, bit11Carry)
		reg.SetFlag(CFlag, bit15Carry)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
