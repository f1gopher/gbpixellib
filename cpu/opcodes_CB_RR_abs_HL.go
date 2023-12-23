package cpu

import (
	"errors"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_RR_abs_HL struct {
	opcodeBase

	value uint8
	bit0  bool
}

func createCB_RR_abs_HL(opcode uint8) *opcode_CB_RR_abs_HL {
	return &opcode_CB_RR_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RR (HL)",
			opcodeLength: 2, // TODO - check
		},
	}
}

func (o *opcode_CB_RR_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		o.bit0 = memory.GetBit(o.value, 0)
		o.value = o.value >> 1

		// TODO - no idea if this is right
		if reg.GetFlag(CFlag) {
			o.value ^= 0x80
		}
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
