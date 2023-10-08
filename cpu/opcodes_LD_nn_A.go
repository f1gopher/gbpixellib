package cpu

import "errors"

type opcode_LD_nn_A struct {
	opcodeBase
}

func createLD_nn_A(opcode uint8) *opcode_LD_nn_A {
	return &opcode_LD_nn_A{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD (nn),A",
			opcodeLength: 3,
		},
	}
}

func (o *opcode_LD_nn_A) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		lsb := readAndIncPC(reg, mem)
		msb := readAndIncPC(reg, mem)
		nn := combineBytes(msb, lsb)
		mem.WriteByte(nn, reg.Get8(A))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
