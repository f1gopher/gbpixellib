package cpu

import (
	"errors"
	"fmt"
)

type opcode_CB_SWAP_r struct {
	opcodeBase

	target Register
}

func createCB_SWAP_r(opcode uint8, reg Register) *opcode_CB_SWAP_r {
	return &opcode_CB_SWAP_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SWAP %s", reg.String()),
			opcodeLength: 2, // TODO - check
		},
		target: reg,
	}
}

func (o *opcode_CB_SWAP_r) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		value := reg.Get8(o.target)
		upper := (value & 0xF0) >> 4
		lower := (value & 0x0F) << 4
		value = upper ^ lower
		reg.Set8(o.target, value)

		reg.SetFlag(ZFlag, value == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, false)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
