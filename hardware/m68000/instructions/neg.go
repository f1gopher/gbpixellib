package instructions

import "fmt"

type neg struct {
	length int
	size   instSize
	src    addressModeInfo
}

func CreateNEG(value uint16, code []uint8, codeOffset int) *neg {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	return &neg{
		length: 2 + instSize,
		size:   size,
		src:    src,
	}
}

func (n *neg) Name() string {
	return fmt.Sprintf("neg.%s   %s", n.size.string(), n.src.string())
}

func (n *neg) Length() int {
	return n.length
}
