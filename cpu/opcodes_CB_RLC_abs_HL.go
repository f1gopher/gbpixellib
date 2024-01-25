package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcodeCB_RLC_abs_HL struct {
	opcodeBase

	data uint8
	bit7 bool
}

func createCB_RLC_abs_HL(opcode uint8) *opcodeCB_RLC_abs_HL {
	return &opcodeCB_RLC_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RLC (HL)"),
			opcodeLength: 1,
		},
	}
}

func (o *opcodeCB_RLC_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.data = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		o.bit7 = memory.GetBit(o.data, 7)
		o.data = o.data << 1

		if o.bit7 {
			o.data ^= 0x01
		}
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(HL), o.data)

		reg.SetFlag(ZFlag, o.data == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, o.bit7)
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invlaid cycle")
}
