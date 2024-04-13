package cpu

import (
	"errors"
	"fmt"
)

type opcode_RET_cc struct {
	opcodeBase

	flag      RegisterFlags
	modifier  bool
	condition bool
	msb       uint8
	lsb       uint8
}

func createRET_cc(opcode uint8, flag RegisterFlags, modifier bool) *opcode_RET_cc {
	flagText := flag.String()
	if !modifier {
		flagText = "N" + flagText
	}

	return &opcode_RET_cc{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("RET %s", flagText),
			opcodeLength: 1,
		},
		flag:     flag,
		modifier: modifier,
	}
}

func (o *opcode_RET_cc) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		// Do NC or NZ by settings the modifier to false
		if o.modifier {
			o.condition = reg.GetFlag(o.flag)
		} else {
			o.condition = !reg.GetFlag(o.flag)
		}
		return false, nil
	}

	if cycleNumber == 2 {
		return !o.condition, nil
	}

	if !o.condition {
		return false, errors.New("Invalid cycle for condition")
	}

	if cycleNumber == 3 {
		o.lsb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 4 {
		o.msb = readAndIncSP(reg, mem)
		return false, nil
	}

	if cycleNumber == 5 {
		reg.Set16(PC, CombineBytes(o.msb, o.lsb))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
