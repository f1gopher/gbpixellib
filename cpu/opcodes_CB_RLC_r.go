package cpu

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type opcodeCB_RLC_r struct {
	opcodeBase

	src Register
}

func createCB_RLC_r(opcode uint8, reg Register) *opcodeCB_RLC_r {
	return &opcodeCB_RLC_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RLC %s", reg.String()),
			opcodeLength: 1,
		},
		src: reg,
	}
}

func (o *opcodeCB_RLC_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		value := reg.Get8(o.src)
		bit7 := memory.GetBit(value, 7)
		result := value << 1

		if bit7 {
			result ^= 0x01
		}

		reg.Set8(o.src, result)

		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, bit7)
		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invlaid cycle")
}
