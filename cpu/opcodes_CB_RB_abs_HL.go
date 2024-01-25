package cpu

import (
	"errors"
)

type opcode_CB_RL_abs_HL struct {
	opcodeBase

	value uint8
	carry uint8
}

func createCB_RL_abs_HL(opcode uint8) *opcode_CB_RL_abs_HL {
	return &opcode_CB_RL_abs_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RL (HL)",
			opcodeLength: 2, // TODO - check this
		},
	}
}

func (o *opcode_CB_RL_abs_HL) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = mem.ReadByte(reg.Get16(HL))
		return false, nil
	}

	if cycleNumber == 2 {
		o.carry = o.value & 0x80
		o.value = o.value << 1
		if reg.GetFlag(CFlag) {
			o.value = o.value ^ 0x01
		}
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(HL), o.value)

		reg.SetFlag(ZFlag, o.value == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, o.carry == 0x80)

		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
