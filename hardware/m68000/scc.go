package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type scc struct {
	length    int
	dest      addressModeInfo
	condition testCondition
}

func CreateScc(value uint16, code []uint8, codeOffset int) *scc {
	dest, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	return &scc{
		length:    2 + instSize,
		dest:      dest,
		condition: condition(value),
	}
}

func (s *scc) Name() string {
	return fmt.Sprintf(
		"s%-3s    %s",
		s.condition.string(),
		s.dest.string())
}

func (s *scc) Length() int {
	return s.length
}

func (s *scc) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
