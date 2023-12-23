package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_SRL_r struct {
	opcodeBase

	target Register
}

func createCB_SRL_r(opcode uint8, reg Register) *opcode_CB_SRL_r {
	return &opcode_CB_SRL_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SRL %s", reg.String()),
			opcodeLength: 2, // TODO - check
		},
		target: reg,
	}
}

func (o *opcode_CB_SRL_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		value := reg.Get8(o.target)
		bit0 := memory.GetBit(value, 0)
		value = value >> 1
		reg.Set8(o.target, value)

		reg.SetFlag(ZFlag, value == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, bit0)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
