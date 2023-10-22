package cpu

import (
	"errors"
)

type opcode_LD_HL_SP_plus_e struct {
	opcodeBase

	n int8
}

func createLD_HL_SP_plus_e(opcode uint8) *opcode_LD_HL_SP_plus_e {
	return &opcode_LD_HL_SP_plus_e{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD HL,SP+e",
			opcodeLength: 2,
		},
	}
}

func (o *opcode_LD_HL_SP_plus_e) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = int8(readAndIncPC(reg, mem))
		return false, nil
	}

	if cycleNumber == 2 {
		var hl uint16 = 0
		sp := reg.Get16(SP)
		if o.n > 127 {
			hl = sp + uint16(-o.n)
		} else {
			hl = sp + uint16(o.n)
		}
		reg.Set16(HL, hl)

		c := sp ^ uint16(o.n) ^ ((sp + uint16(o.n)) & 0xFFFF)
		cFlag := c&0x0100 == 0x0100
		hFlag := c&0x0010 == 0x0010

		reg.SetFlag(ZFlag, false)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, hFlag)
		reg.SetFlag(CFlag, cFlag)

		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
