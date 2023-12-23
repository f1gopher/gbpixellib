package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_BIT_b_HL struct {
	opcodeBase

	bit   uint8
	hl    uint16
	value uint8
}

func createCB_BIT_b_HL(opcode uint8, bit uint8) *opcode_CB_BIT_b_HL {
	return &opcode_CB_BIT_b_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("BIT %d,(HL)", bit),
			opcodeLength: 2, // TODO - check
		},
		bit: bit,
	}
}

func (o *opcode_CB_BIT_b_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.hl = reg.Get16(HL)
		return false, nil
	}

	if cycleNumber == 2 {
		o.value = mem.ReadByte(o.hl)
		return false, nil
	}

	if cycleNumber == 3 {
		result := memory.GetBit(o.value, int(o.bit))

		reg.SetFlag(ZFlag, result == false)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, true)
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
