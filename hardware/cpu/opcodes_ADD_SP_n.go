package cpu

import (
	"errors"
)

type opcode_ADD_SP_n struct {
	opcodeBase

	n  uint8
	sp uint16
}

func createADD_SP_n(opcode uint8) *opcode_ADD_SP_n {
	return &opcode_ADD_SP_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "ADD SP,n",
			opcodeLength: 4,
		},
	}
}

func (o *opcode_ADD_SP_n) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.n = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.sp = reg.Get16(SP)
		return false, nil
	}

	if cycleNumber == 3 {
		var result uint16
		if o.n > 127 {
			result = o.sp - uint16(-o.n)
		} else {
			result = o.sp + uint16(o.n)
		}

		reg.Set16(SP, result)

		c := o.sp ^ uint16(o.n) ^ ((o.sp + uint16(o.n)) & 0xFFFF)

		reg.SetFlag(ZFlag, false)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, (c&0x10) == 0x10)
		reg.SetFlag(CFlag, (c&0x100) == 0x100)
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
