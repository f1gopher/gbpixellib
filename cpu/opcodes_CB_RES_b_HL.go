package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_RES_b_HL struct {
	opcodeBase

	bit    uint8
	target Register
	hl     uint16
	value  uint8
}

func createCB_RES_b_HL(opcode uint8, bit uint8) *opcode_CB_RES_b_HL {
	return &opcode_CB_RES_b_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RES %d,(HL)", bit),
			opcodeLength: 2, // TODO - Check this
		},
		bit: bit,
	}
}

func (o *opcode_CB_RES_b_HL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.hl = reg.Get16(HL)
		return false, nil
	}

	if cycleNumber == 2 {
		o.value = mem.ReadByte(o.hl)
		return false, nil
	}

	if cycleNumber == 3 {
		o.value = memory.SetBit(o.value, o.bit, false)
		mem.WriteByte(o.hl, o.value)
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
