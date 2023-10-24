package cpu

import (
	"errors"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcode_RLCA struct {
	opcodeBase

	src  Register
	dest Register
}

func createRLCA(opcode uint8) *opcode_RLCA {
	return &opcode_RLCA{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RLCA",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_RLCA) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		return false, errors.New("Invalid cycle")
	}

	value := reg.Get8(A)
	bit7 := memory.GetBit(value, 7)
	result := value << 1
	reg.Set8(A, result)

	reg.SetFlag(ZFlag, result == 0)
	reg.SetFlag(NFlag, false)
	reg.SetFlag(HFlag, false)
	reg.SetFlag(CFlag, bit7)

	return true, nil
}
