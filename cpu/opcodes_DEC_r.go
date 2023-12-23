package cpu

import (
	"fmt"
)

type opcode_DEC_r struct {
	opcodeBase

	target Register
}

func createDEC_r(opcode uint8, reg Register) *opcode_DEC_r {
	return &opcode_DEC_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("DEC %s", reg.String()),
			opcodeLength: 1,
		},
		target: reg,
	}
}

func (o *opcode_DEC_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		panic("Invalid cycle")
	}

	original := reg.Get8(o.target)

	result, carryBit3, _ := subtract8BitWithCarry(original, 1)

	reg.Set8(o.target, result)
	reg.SetFlag(ZFlag, result == 0)
	reg.SetFlag(NFlag, true)
	reg.SetFlag(HFlag, carryBit3)

	return true, nil
}
