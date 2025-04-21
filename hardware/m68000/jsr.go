package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type jsr struct {
	instSize int
	src      addressModeInfo

	data uint32
}

func CreateJSR(value uint16, code []uint8, codeOffset int) *jsr {
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	return &jsr{
		instSize: 2 + instSize,
		src:      src,
	}
}

func (j *jsr) Name() string {
	return fmt.Sprintf(
		"jsr     %s",
		j.src.string())
}

func (j *jsr) Length() int {
	return j.instSize
}

func (j *jsr) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
