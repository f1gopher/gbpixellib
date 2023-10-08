package cpu

import (
	"errors"
)

type opcode_LD_A_HL_dec struct {
	opcodeBase

	src  register
	dest register
}

func createLD_A_HL_dec(opcode uint8) *opcode_LD_A_HL_dec {
	return &opcode_LD_A_HL_dec{
		opcodeBase: opcodeBase{
			opcodeId:   opcode,
			opcodeName: "LD A,(HL-)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_LD_A_HL_dec) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
	reg.set8(A, mem.ReadByte(reg.Get16(HL)))	
		reg.set16(HL, reg.Get16(HL)-1)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
