package cpu

import (
	"fmt"
)

type opcode_ADD_rr_rr struct {
	opcodeBase

	src  register
	dest register
}

func createADD_rr_rr(opcode uint8, srcReg register, destReg register) *opcode_ADD_rr_rr {
	return &opcode_ADD_rr_rr{
		opcodeBase: opcodeBase{
			opcodeId:     opcode,
			opcodeName:   fmt.Sprintf("ADD %s,%s", srcReg.String(), destReg.String()),
			opcodeLength: 1, // TODO - check
		},
		src:  srcReg,
		dest: destReg,
	}
}

func (o *opcode_ADD_rr_rr) doCycle(cycleNumber int, reg registersInterface, mem memoryInterface) (completed bool, err error) {

	panic("Not implemented - ADD_rr_rr")

	return true, nil
}
