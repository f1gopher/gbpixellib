package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type sub struct {
	length  int
	size    instSize
	src     addressModeInfo
	mode    byte
	destReg byte
}

func CreateSUB(value uint16, code []uint8, codeOffset int) *sub {
	mode := valueFrom3Bits(value, 6)
	size := byteSize
	if mode == 0b001 || mode == 0b101 {
		size = wordSize
	} else if mode == 0b010 || mode == 0b110 {
		size = longSize
	}
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	destReg := valueFrom3Bits(value, 9)

	return &sub{
		length:  2 + instSize,
		size:    size,
		src:     src,
		mode:    mode,
		destReg: destReg,
	}
}

func (s *sub) Name() string {
	if s.mode == 0b0 || s.mode == 0b00000001 || s.mode == 0b00000010 {
		return fmt.Sprintf(
			"sub.%s   %s, D%d",
			s.size.string(),
			s.src.string(),
			s.destReg)
	}

	return fmt.Sprintf(
		"sub.%s   D%d, %s",
		s.size.string(),
		s.destReg,
		s.src.string())
}

func (s *sub) Length() int {
	return s.length
}

func (s *sub) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
