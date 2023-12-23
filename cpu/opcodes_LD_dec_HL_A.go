package cpu

import (
	"errors"
)

type opcode_LD_dec_HL_A struct {
	opcodeBase

	src  Register
	dest Register
}

func createLD_dec_HL_A(opcode uint8) *opcode_LD_dec_HL_A {
	return &opcode_LD_dec_HL_A{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD (HL-),A",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_LD_dec_HL_A) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		mem.WriteByte(reg.Get16(HL), reg.Get8(A))
		reg.Set16(HL, reg.Get16(HL)-1)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
