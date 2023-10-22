package cpu

import "errors"

type opcode_CCF struct {
	opcodeBase
}

func createCCF(opcode uint8) *opcode_CCF {
	return &opcode_CCF{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "CCF",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_CCF) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.SetFlag(NFlag, false)
		reg.SetFlag(HFlag, false)
		reg.SetFlag(CFlag, !reg.GetFlag(CFlag))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
