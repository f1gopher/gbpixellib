package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type oriToCCR struct {
	data byte
}

func CreateORIToCCR(value uint16, code []uint8, codeOffset int) *oriToCCR {
	return &oriToCCR{
		data: byte(readData(code, codeOffset, byteSize)),
	}
}

func (o *oriToCCR) Name() string {
	return fmt.Sprintf("ori     #$%x, CCR", o.data)
}

func (o *oriToCCR) Length() int {
	return 4
}

func (o *oriToCCR) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
