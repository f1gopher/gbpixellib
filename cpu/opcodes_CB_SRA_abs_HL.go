package cpu

import (
	"errors"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcodeCB_SRA_abs_HL struct {
	opcodeBase

	value uint8
	bit0  bool
}

func createCB_SRA_abs_HL(opcode uint8) *opcodeCB_SRA_abs_HL {
	return &opcodeCB_SRA_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "SRA (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcodeCB_SRA_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		o.bit0 = memory.GetBit(o.value, 0)
		// no idea
		o.value = (o.value >> 1) | (o.value & 0x80)

		//	// TODO - is this right?
		//	if reg.GetFlag(CFlag) {
		//		o.value ^= 0x80
		//	}
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(HL), o.value)

		reg.SetFlag(ZFlag, o.value == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, o.bit0)
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
