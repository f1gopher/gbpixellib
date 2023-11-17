package cpu

import (
	"errors"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_RRC_abs_HL struct {
	opcodeBase

	data uint8
	bit0 bool
}

func createCB_RRC_abs_HL(opcode uint8) *opcode_CB_RRC_abs_HL {
	return &opcode_CB_RRC_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RRC (HL)",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_CB_RRC_abs_HL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.data = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		o.bit0 = memory.GetBit(o.data, 0)
		o.data = o.data >> 1

		// TODO - what the hell and why???
		if o.bit0 {
			o.data ^= 0x80
		}
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(HL), o.data)

		reg.SetFlag(ZFlag, o.data == 0)
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
