package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type suba struct {
	length int
	size   instSize
	src    addressModeInfo
	reg    byte
}

func CreateSUBA(value uint16, code []uint8, codeOffset int) *suba {
	opmode := valueFrom3Bits(value, 6)
	size := wordSize
	if opmode == 0b00000111 {
		size = longSize
	}
	reg := valueFrom3Bits(value, 9)
	src, _, instSize := readEffectiveAddressMode(value, 0, zeroSize, code, codeOffset, codeOffset)

	if src.mode == invalidAddressMode {
		return nil
	}

	return &suba{
		length: 2 + instSize,
		size:   size,
		src:    src,
		reg:    reg,
	}
}

func (s *suba) Name() string {
	return fmt.Sprintf(
		"suba.%s  %s, A%d",
		s.size.string(),
		s.src.string(),
		s.reg)
}

func (s *suba) Length() int {
	return s.length
}

func (s *suba) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
