package cpu

import (
	"errors"
)

type opcode_JP_nn struct {
	opcodeBase

	msb uint8
	lsb uint8
}

func createJP_nn(opcode uint8) *opcode_JP_nn {
	return &opcode_JP_nn{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "JP nn",
			opcodeLength: 3,
		},
	}
}

func (o *opcode_JP_nn) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		reg.Set16(PC, CombineBytes(o.msb, o.lsb))
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
