package cpu

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_BIT_b_r struct {
	opcodeBase

	bit    uint8
	target register
}

func createCB_BIT_b_r(opcode uint8, bit uint8, reg register) *opcode_CB_BIT_b_r {
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

func (o *opcode_CB_BIT_b_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		panic("Invalid cycle")
	}

	value := reg.Get8(o.target)
	result := memory.GetBit(value, int(o.bit))

	reg.setFlag(ZFlag, result == false)
	reg.setFlag(NFlag, false)
	reg.setFlag(HFlag, true)

	return true, nil
}
