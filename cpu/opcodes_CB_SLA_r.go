package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_SLA_r struct {
	opcodeBase

	target Register
}

func createCB_SLA_r(opcode uint8, reg Register) *opcode_CB_SLA_r {
	return &opcode_CB_SLA_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SLA %s", reg.String()),
			opcodeLength: 2, // TODO - check this
		},
		target: reg,
	}
}

func (o *opcode_CB_SLA_r) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		value := reg.Get8(o.target)
		bit7 := memory.GetBit(value, 7)
		value = value << 1
		reg.Set8(o.target, value)

		reg.SetFlag(ZFlag, value == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, bit7)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
