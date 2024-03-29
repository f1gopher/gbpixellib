package cpu

import (
	"errors"
)

type opcode_RETI struct {
	opcodeBase

	msb uint8
	lsb uint8
}

func createRETI(opcode uint8) *opcode_RETI {
	return &opcode_RETI{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   "RETI",
			opcodeLength: 1,
		},
	}
}

func (o *opcode_RETI) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		reg.Set16(PC, CombineBytes(o.msb, o.lsb))
		return false, nil
	}

	if cycleNumber == 4 {
		//mem.WriteByte(0xFFFF, 0x01)
		reg.SetIME(true)
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
