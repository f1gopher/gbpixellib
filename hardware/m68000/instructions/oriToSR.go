package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type oriToSR struct {
	data uint16
}

func CreateORIToSR(value uint16, code []uint8, codeOffset int) *oriToSR {
	return &oriToSR{
		data: uint16(readData(code, codeOffset, wordSize)),
	}
}

func (o *oriToSR) Name() string {
	return fmt.Sprintf("ori     #$%x, SR", o.data)
}

func (o *oriToSR) Length() int {
	return 4
}

func (o *oriToSR) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
