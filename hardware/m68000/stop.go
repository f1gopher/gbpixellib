package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type stop struct {
	immediate uint16
}

func CreateSTOP(value uint16, code []uint8, codeOffset int) *stop {
	return &stop{
		immediate: uint16(readData(code, codeOffset, wordSize)),
	}
}

func (s *stop) Name() string {
	return fmt.Sprintf("stop    #$%x", s.immediate)
}

func (s *stop) Length() int {
	return 4
}

func (s *stop) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
