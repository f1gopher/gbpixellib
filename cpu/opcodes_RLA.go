package cpu

import (
	"errors"
)

type opcode_RLA struct {
	opcodeBase

	src  register
	dest register
}

func createRLA(opcode uint8) *opcode_RLA {
	return &opcode_RLA{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RLA",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_RLA) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber != 1 {
		return false, errors.New("Invalid cycle")
	}

	value := reg.Get8(A)
	carry := value & 0x80
	result := value << 1
	if reg.GetFlag(CFlag) {
		result = result ^ 0x01
	}
	reg.set8(A, result)

	reg.setFlag(ZFlag, result == 0)
	reg.setFlag(NFlag, false)
	reg.setFlag(HFlag, false)
	reg.setFlag(CFlag, carry == 0x80)

	return true, nil
}
