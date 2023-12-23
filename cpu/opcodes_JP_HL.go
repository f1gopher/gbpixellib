package cpu

import "errors"

type opcode_JP_HL struct {
	opcodeBase
}

func createJP_HL(opcode uint8) *opcode_JP_HL {
	return &opcode_JP_HL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "JP HL",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_JP_HL) doCycle(cycleNumber int, reg RegistersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		reg.Set16(PC, reg.Get16(HL))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
