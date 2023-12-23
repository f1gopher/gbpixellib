package cpu

import (
	"errors"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_RRCA struct {
	opcodeBase
}

func createRRCA(opcode uint8) *opcode_RRCA {
	return &opcode_RRCA{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RRCA",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_RRCA) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		a := reg.Get8(A)
		bit0 := memory.GetBit(a, 0)
		a = a >> 1

		// TODO - what the hell and why???
		if bit0 {
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
