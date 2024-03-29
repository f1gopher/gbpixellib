package cpu

import (
	"errors"
	"fmt"
)

type opcode_JR_cc_e struct {
	opcodeBase

	flag      RegisterFlags
	modifier  bool
	condition bool
	e         int8
}

func createJR_cc_e(opcode uint8, flag RegisterFlags, modifier bool) *opcode_JR_cc_e {
	flagText := flag.String()
	if !modifier {
		flagText = "N" + flagText
	}

	return &opcode_JR_cc_e{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("JR %s,e", flagText),
			opcodeLength: 2,
		},
		flag:     flag,
		modifier: modifier,
	}
}

func (o *opcode_JR_cc_e) doCycle(cycleNumber int, reg RegistersInterface, mem MemoryInterface) (completed bool, err error) {

	if cycleNumber == 1 {
		o.e = int8(readAndIncPC(reg, mem))
		return false, nil
	}

	if cycleNumber == 2 {
		o.condition = reg.GetFlag(o.flag)
		return o.condition != o.modifier, nil
	}

	if o.condition != o.modifier {
		return false, errors.New("Invalid cycle for condition")
	}

	if cycleNumber == 3 {
		reg.Set16(PC, adds8Tou16(reg.Get16(PC), o.e))
		return true, nil
	}

	return false, errors.New("Invalid cycle")
}
