package instructions

import (
	"fmt"
)

type not struct {
	length int
	size   instSize
	src    addressModeInfo
}

func CreateNOT(value uint16, code []uint8, codeOffset int) *not {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	return &not{
		length: 2 + instSize,
		size:   size,
		src:    src,
	}
}

func (n *not) Name() string {
	return fmt.Sprintf(
		"not.%s   %s",
		n.size.string(),
		n.src.string())
}

func (n *not) Length() int {
	return n.length
}
