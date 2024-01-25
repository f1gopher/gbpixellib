package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcodeCB_SRA_r struct {
	opcodeBase

	src Register
}

func createCB_SRA_r(opcode uint8, reg Register) *opcodeCB_SRA_r {
	return &opcodeCB_SRA_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SRA %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcodeCB_SRA_r) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		a := reg.Get8(o.src)
		bit0 := memory.GetBit(a, 0)
		// no idea
		a = (a >> 1) | (a & 0x80)

		//		// TODO - is this right?
		//		if reg.GetFlag(CFlag) {
		//			a ^= 0x80
		//		}

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
