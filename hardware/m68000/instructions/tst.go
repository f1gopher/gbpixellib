package instructions

import "fmt"

type tst struct {
	length int
	size   instSize
	src    addressModeInfo
}

func CreateTST(value uint16, code []uint8, codeOffset int) *tst {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
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
