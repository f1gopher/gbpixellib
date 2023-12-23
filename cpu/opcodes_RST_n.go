package cpu

import (
	"errors"
	"fmt"
)

type opcode_RST_n struct {
	opcodeBase

	value uint8
}

func createRST_n(opcode uint8, value uint8) *opcode_RST_n {
	return &opcode_RST_n{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RST X%02x", value),
			opcodeLength: 1,
		},
		value: value,
	}
}

func (o *opcode_RST_n) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		decSP(reg)
		return false, nil
	}

	if cycleNumber == 2 {
		mem.WriteByte(reg.Get16(SP), Msb(reg.Get16(PC)))
		decSP(reg)
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(reg.Get16(SP), Lsb(reg.Get16(PC)))
		return false, nil
	}

	if cycleNumber == 4 {
		// TODO - is o.value right?
		reg.Set16(PC, combineBytes(0x00, o.value))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
