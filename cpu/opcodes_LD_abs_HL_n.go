package cpu

import (
	"errors"
)

type opcode_LD_abs_HL_n struct {
	opcodeBase

	value uint8
}

func createLD_abs_HL_n(opcode uint8) *opcode_LD_abs_HL_n {
	return &opcode_LD_abs_HL_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD (HL),n",
			opcodeLength: 2,
		},
	}
}

func (o *opcode_LD_abs_HL_n) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.value = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		mem.WriteByte(reg.Get16(HL), o.value)
		return false, nil
	}

	if cycleNumber == 3 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
