package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type tst struct {
	length int
	size   instSize
	src    addressModeInfo
}

func CreateTST(value uint16, code []uint8, codeOffset int) *tst {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	if src.mode == addressRegisterDirect ||
		src.mode == programCounterIndirectWithDisplacement ||
		src.mode == preIndexedImediate ||
		src.mode == programCounterIndirectWithIndex8bitDisplacement {
		return nil
	}

	return &tst{
		length: 2 + instSize,
		size:   size,
		src:    src,
	}
}

func (t *tst) Name() string {
	return fmt.Sprintf(
		"tst.%s   %s",
		t.size.string(),
		t.src.string())
}

func (t *tst) Length() int {
	return t.length
}

func (t *tst) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
