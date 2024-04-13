package cpu

import "errors"

type opcode_LD_nn_A struct {
	opcodeBase

	msb uint8
	lsb uint8
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

func (o *opcode_LD_nn_A) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		nn := CombineBytes(o.msb, o.lsb)
		mem.WriteByte(nn, reg.Get8(A))
		return false, nil
	}

	if cycleNumber == 4 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
