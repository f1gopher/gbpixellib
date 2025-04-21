package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type jmp struct {
	length int
	src    addressModeInfo
}

func CreateJMP(value uint16, code []uint8, codeOffset int) *jmp {
	src, _, instSize := readEffectiveAddressMode(value, 0, zeroSize, code, codeOffset, codeOffset)
	return &jmp{
		length: 2 + instSize,
		src:    src,
	}
}

func (j *jmp) Name() string {
	return fmt.Sprintf(
		"jmp     %s",
		j.src.string())
}

func (j *jmp) Length() int {
	return j.length
}

func (j *jmp) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
