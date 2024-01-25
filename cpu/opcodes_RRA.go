package cpu

import (
	"errors"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_RRA struct {
	opcodeBase
}

func createRRA(opcode uint8) *opcode_RRA {
	return &opcode_RRA{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RRA",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_RRA) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		a := reg.Get8(A)
		bit0 := memory.GetBit(a, 0)
		a = a >> 1

		// TODO - is this right?
		if reg.GetFlag(CFlag) {
			a ^= 0x80
		}

		reg.Set8(A, a)

		reg.SetFlag(ZFlag, a == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, bit0)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
