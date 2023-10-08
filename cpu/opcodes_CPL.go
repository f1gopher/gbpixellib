package cpu

import "errors"

type opcode_CPL struct {
	opcodeBase
}

func createCPL(opcode uint8) *opcode_CPL {
	return &opcode_CPL{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "CPL",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_CPL) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		a := reg.Get8(A)
		a = a ^ a
		reg.set8(A, a)
		reg.setFlag(NFlag, true)
		reg.setFlag(HFlag, true)
	}

	return false, errors.New("Invalid cycle")
}