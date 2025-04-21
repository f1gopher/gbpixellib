package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type mulu struct {
	length int
	src    addressModeInfo
	reg    byte
}

func CreateMULU(value uint16, code []uint8, codeOffset int) *mulu {
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	reg := valueFrom3Bits(value, 9)

	return &mulu{
		length: 2 + instSize,
		src:    src,
		reg:    reg,
	}
}

func (m *mulu) Name() string {
	return fmt.Sprintf(
		"mulu.w  %s, D%d",
		m.src.string(),
		m.reg)
}

func (m *mulu) Length() int {
	return m.length
}

func (m *mulu) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
