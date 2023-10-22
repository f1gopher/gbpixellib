package cpu

import "errors"

type opcode_LD_A_nn struct {
	opcodeBase

	msb uint8
	lsb uint8
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
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		nn := combineBytes(o.msb, o.lsb)
		reg.Set8(A, mem.ReadByte(nn))
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
