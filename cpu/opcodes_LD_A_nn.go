package cpu

import "errors"

type opcode_LD_A_nn struct {
	opcodeBase
}

func createLD_A_nn(opcode uint8) *opcode_LD_A_nn {
	return &opcode_LD_A_nn{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD A,(nn)",
			opcodeLength: 3,
		},
	}
}

func (o *opcode_LD_A_nn) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		lsb := readAndIncPC(reg, mem)
		msb := readAndIncPC(reg, mem)
		nn := combineBytes(msb, lsb)
		reg.set8(A, mem.ReadByte(nn))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
