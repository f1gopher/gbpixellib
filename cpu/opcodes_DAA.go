package cpu

import (
	"errors"
)

type opcode_DAA struct {
	opcodeBase
}

func createDAA(opcode uint8) *opcode_DAA {
	return &opcode_DAA{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "DAA",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_DAA) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {

		value := uint16(reg.Get8(A))

		if reg.GetFlag(NFlag) == false {
			if reg.GetFlag(HFlag) || value&0x0F > 9 {
				value += 0x06
			}
			if reg.GetFlag(CFlag) || value > 0x9F {
				value += 0x60
			}
		} else {
			if reg.GetFlag(HFlag) {
				value = (value - 0x06) & 0xFF
			}
			if reg.GetFlag(CFlag) {
				value -= 0x60
			}
		}

		if value&0x100 == 0x100 {
			reg.SetFlag(CFlag, true)
		}
		value &= 0xFF
		reg.Set8(A, uint8(value))

		reg.SetFlag(ZFlag, value == 0)
		reg.SetFlag(HFlag, false)

		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
