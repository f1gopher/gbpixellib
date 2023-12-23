package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_BIT_b_r struct {
	opcodeBase

	bit    uint8
	target Register
}

func createCB_BIT_b_r(opcode uint8, bit uint8, reg Register) *opcode_CB_BIT_b_r {
	return &opcode_CB_BIT_b_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("BIT %d,%s", bit, reg.String()),
			opcodeLength: 2, // TODO - check
		},
		bit:    bit,
		target: reg,
	}
}

func (o *opcode_CB_BIT_b_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		value := reg.Get8(o.target)
		result := memory.GetBit(value, int(o.bit))

		reg.SetFlag(ZFlag, result == false)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, true)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
