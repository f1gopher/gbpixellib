package cpu

import "errors"

type opcode_SCF struct {
	opcodeBase
}

func createSCF(opcode uint8) *opcode_SCF {
	return &opcode_SCF{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "SCF",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_SCF) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, true)
	}

	return false, errors.New("Invalid cycle")
}
