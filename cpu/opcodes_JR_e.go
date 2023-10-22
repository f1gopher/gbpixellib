package cpu

import (
	"errors"
)

type opcode_JR_e struct {
	opcodeBase

	e int8
}

func createJR_e(opcode uint8) *opcode_JR_e {
	return &opcode_JR_e{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "JR e",
			opcodeLength: 2,
		},
	}
}

func (o *opcode_JR_e) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.e = int8(readAndIncPC(reg, mem))
		return false, nil
	}

	if cycleNumber == 2 {
		pc := reg.Get16(PC)
		// TODO - better unit test this
		if o.e > 0 {
			pc = pc + uint16(o.e)
		} else {
			pc = pc - uint16(-o.e)
		}
		reg.Set16(PC, pc)
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
