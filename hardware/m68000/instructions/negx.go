package instructions

import "fmt"

type negx struct {
	length int
	size   instSize
	src    addressModeInfo
}

func CreateNEGX(value uint16, code []uint8, codeOffset int) *negx {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	return &negx{
		length: 2 + instSize,
		size:   size,
		src:    src,
	}
}

func (n *negx) Name() string {
	return fmt.Sprintf(
		"negx.%s  %s",
		n.size.string(),
		n.src.string())
}

func (n *negx) Length() int {
	return n.length
}
