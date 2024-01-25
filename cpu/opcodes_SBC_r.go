package cpu

import (
	"errors"
	"fmt"
)

type opcode_SBC_r struct {
	opcodeBase

	src Register
}

func createSBC_r(opcode uint8, reg Register) *opcode_SBC_r {
	return &opcode_SBC_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("SBC %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcode_SBC_r) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		result, bit4Carry, noBorrow := subtract8BitWithCarryBit4(reg.Get8(A), reg.Get8(o.src), reg.GetFlag(CFlag))
		reg.Set8(A, result)
		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, true)
		reg.SetFlag(HFlag, bit4Carry)
		reg.SetFlag(CFlag, noBorrow)

		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
