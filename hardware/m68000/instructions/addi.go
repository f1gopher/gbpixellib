package instructions

import (
	"fmt"
)

type addi struct {
	length int
	opSize instSize
	dest   addressModeInfo
	data   uint32
}

func CreateADDI(value uint16, code []uint8, codeOffset int) *addi {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset+size.bytesUsed(), codeOffset)
	data := readData(code, codeOffset, size)

	return &addi{
		length: 2 + instSize + size.bytesUsed(),
		opSize: size,
		dest:   dest,

		data: data,
	}
}

func (a *addi) Name() string {
	return fmt.Sprintf(
		"addi.%s  #$%x, %s",
		a.opSize.string(),
		a.data,
		a.dest.string())
}

func (a *addi) Length() int {
	return a.length
}
