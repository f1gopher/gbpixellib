package cpu

import (
	"errors"
)

type opcode_CB_SWAP_abs_HL struct {
	opcodeBase

	value uint8
}

func createCB_SWAP_abs_HL(opcode uint8) *opcode_CB_SWAP_abs_HL {
	return &opcode_CB_SWAP_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "SWAP (HL)",
			opcodeLength: 2, // TODO - check
		},
	}
}

func (o *opcode_CB_SWAP_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		upper := (o.value & 0xF0) >> 4
		lower := (o.value & 0x0F) << 4
		o.value = upper ^ lower
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(HL), o.value)

		reg.SetFlag(ZFlag, o.value == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, false)
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
