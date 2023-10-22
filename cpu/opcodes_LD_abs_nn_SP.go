package cpu

import (
	"errors"
)

type opcode_LD_abs_nn_SP struct {
	opcodeBase

	msb uint8
	lsb uint8
	nn  uint16
}

func createLD_abs_nn_SP(opcode uint8) *opcode_LD_abs_nn_SP {
	return &opcode_LD_abs_nn_SP{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "LD (nn),SP",
			opcodeLength: 3,
		},
	}
}

func (o *opcode_LD_abs_nn_SP) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncPC(reg, mem)
		o.nn = combineBytes(o.msb, o.lsb)
		return false, nil
	}

	if cycleNumber == 3 {
		mem.WriteByte(o.nn, Lsb(reg.Get16(SP)))
		return false, nil
	}

	if cycleNumber == 4 {
		mem.WriteByte(o.nn+1, Msb(reg.Get16(SP)))
		return false, nil
	}

	if cycleNumber == 5 {
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
