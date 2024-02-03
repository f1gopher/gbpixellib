package cpu

import (
	"errors"
)

type opcode_CALL_nn struct {
	opcodeBase

	msb uint8
	lsb uint8
}

func createCALL_nn(opcode uint8) *opcode_CALL_nn {
	return &opcode_CALL_nn{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "CALL nn",
			opcodeLength: 3,
		},
	}
}

func (o *opcode_CALL_nn) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		decSP(reg)
		return false, nil
	}

	if cycleNumber == 4 {
		mem.WriteByte(reg.Get16(SP), Msb(reg.Get16(PC)))
		decSP(reg)
		return false, nil
	}

	if cycleNumber == 5 {
		mem.WriteByte(reg.Get16(SP), Lsb(reg.Get16(PC)))
		return false, nil
	}

	if cycleNumber == 6 {
		reg.Set16(PC, CombineBytes(o.msb, o.lsb))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
