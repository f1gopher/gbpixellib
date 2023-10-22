package cpu

import (
	"errors"
	"fmt"
)

type opcode_JP_cc_nn struct {
	opcodeBase

	flag      registerFlags
	modifier  bool
	condition bool
	msb       uint8
	lsb       uint8
}

func createJP_cc_nn(opcode uint8, flag registerFlags, modifier bool) *opcode_JP_cc_nn {
	flagText := flag.String()
	if !modifier {
		flagText = "N" + flagText
	}

	return &opcode_JP_cc_nn{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("JP %s,nn", flagText),
			opcodeLength: 3,
		},
		flag:     flag,
		modifier: modifier,
	}
}

func (o *opcode_JP_cc_nn) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.lsb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 2 {
		o.msb = readAndIncPC(reg, mem)
		return false, nil
	}

	if cycleNumber == 3 {
		// Do NC or NZ by settings the modifier to false
		if o.modifier {
			o.condition = reg.GetFlag(o.flag)
		} else {
			o.condition = !reg.GetFlag(o.flag)
		}
		return !o.condition, nil
	}

	if !o.condition {
		return false, errors.New("Invalid cycle for condition")
	}

	if cycleNumber == 4 {
		reg.Set16(PC, combineBytes(o.msb, o.lsb))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
