package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type eori struct {
	length int
	size   instSize
	src    addressModeInfo
	dest   addressModeInfo
	data   uint32
}

func CreateEORI(value uint16, code []uint8, codeOffset int) *eori {
	size := twoBitSize(value, 6)
	src := addressModeInfo{
		mode: preIndexedImediate,
		data: readData(code, codeOffset, size),
	}
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+2, codeOffset)

	if dest.mode == invalidAddressMode ||
		dest.mode == addressRegisterDirect ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement {
		return nil
	}

	return &eori{
		length: 2 + instSize + size.bytesUsed(),
		size:   size,
		src:    src,
		dest:   dest,
		data:   readData(code, codeOffset, size),
	}
}

func (e *eori) Name() string {
	return fmt.Sprintf(
		"eori.%s  %s, %s",
		e.size.string(),
		e.src.string(),
		e.dest.string())
}

func (e *eori) Length() int {
	return e.length
}

func (e *eori) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
