package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type move struct {
	length int
	src    addressModeInfo
	dest   addressModeInfo

	size instSize
}

func CreateMOVE(value uint16, code []uint8, codeOffset int) *move {
	size := moveSize(value, 12)
	src, _, srcInstSize := readSrcAddressMode(value, 0, size, code, codeOffset, codeOffset)
	dest, _, destInstSize := readDestAddressMode(value, 6, size, code, codeOffset+srcInstSize, codeOffset)

	if dest.mode == addressRegisterDirect ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement ||
		dest.mode == invalidAddressMode ||
		size == zeroSize ||
		src.mode == invalidAddressMode {
		return nil
	}
	if size == byteSize && src.mode == addressRegisterDirect {
		return nil
	}

	return &move{
		length: 2 + srcInstSize + destInstSize,
		src:    src,
		dest:   dest,
		size:   size,
	}
}

func (m *move) Name() string {
	return fmt.Sprintf(
		"move.%s  %s, %s",
		m.size.string(),
		m.src.string(),
		m.dest.string())
}

func (m *move) Length() int {
	return m.length
}

func (m *move) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
