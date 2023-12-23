package cpu

import (
	"errors"
	"fmt"
)

type opcode_CB_RL_r struct {
	opcodeBase

	target Register
}

func createCB_RL_r(opcode uint8, reg Register) *opcode_CB_RL_r {
	return &opcode_CB_RL_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RL %s", reg.String()),
			opcodeLength: 2, // TODO - check this
		},
		target: reg,
	}
}

func (o *opcode_CB_RL_r) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		value := reg.Get8(o.target)
		carry := value & 0x80
		result := value << 1
		if reg.GetFlag(CFlag) {
			result = result ^ 0x01
		}
		reg.Set8(o.target, result)

		reg.SetFlag(ZFlag, result == 0)
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, carry == 0x80)

		return false, nil
	}

	if cycleNumber == 2 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
