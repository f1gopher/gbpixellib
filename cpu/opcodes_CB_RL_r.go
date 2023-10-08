package cpu

import (
	"fmt"
)

type opcode_CB_RL_r struct {
	opcodeBase

	target register
}

func createCB_RL_r(opcode uint8, reg register) *opcode_CB_RL_r {
	return &opcode_CB_RL_r{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RL %s", reg.String()),
			opcodeLength: 2, // TODO - check this
		},
		target: reg,
	}
}

func (o *opcode_CB_RL_r) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		panic("Invalid cycle")
	}

	value := reg.Get8(o.target)
	carry := value & 0x80
	result := value << 1
	if reg.GetFlag(CFlag) {
		result = result ^ 0x01
	}
	reg.set8(o.target, result)

	reg.setFlag(ZFlag, result == 0)
	reg.setFlag(NFlag, false)
	reg.setFlag(HFlag, false)
	reg.setFlag(CFlag, carry == 0x80)

	return true, nil
}
