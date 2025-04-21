package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type moveToSR struct {
	instSize int
	dest     addressModeInfo
}

func CreateMOVEToSR(value uint16, code []uint8, codeOffset int) *moveToSR {
	dest, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	return &moveToSR{
		instSize: 2 + instSize,
		dest:     dest,
	}
}

func (m *moveToSR) Name() string {
	return fmt.Sprintf(
		"move    %s, SR",
		m.dest.string())
}

func (m *moveToSR) Length() int {
	return m.instSize
}

func (m *moveToSR) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
