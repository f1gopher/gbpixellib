package cpu

import (
	"errors"
)

type opcode_RET struct {
	opcodeBase

	msb uint8
	lsb uint8
}

func createRET(opcode uint8) *opcode_RET {
	return &opcode_RET{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RET",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_RET) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		reg.Set16(PC, combineBytes(o.msb, o.lsb))
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
