package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_CB_RRC_r struct {
	opcodeBase

	src Register
}

func createCB_RRC_r(opcode uint8, reg Register) *opcode_CB_RRC_r {
	return &opcode_CB_RRC_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RRC %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_CB_RRC_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		a := reg.Get8(o.src)
		bit0 := memory.GetBit(a, 0)
		a = a >> 1

		// TODO - what the hell and why???
		if bit0 {
			a ^= 0x80
		}

		reg.Set8(o.src, a)

		reg.SetFlag(ZFlag, a == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, bit0)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
