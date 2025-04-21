package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type andiToCCR struct {
	data byte
}

func CreateANDIToCCR(value uint16, code []uint8, codeOffset int) *andiToCCR {
	return &andiToCCR{
		data: byte(readData(code, codeOffset, byteSize)),
	}
}

func (a *andiToCCR) Name() string {
	return fmt.Sprintf("andi    #$%x, CCR", a.data)
}

func (a *andiToCCR) Length() int {
	return 4
}

func (a *andiToCCR) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
