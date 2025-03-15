package instructions

import "fmt"

type bchg struct {
	length   int
	src      addressModeInfo
	register byte
}

func CreateBCHG(value uint16, code []uint8, codeOffset int) *bchg {
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	register := valueFrom3Bits(value, 9)

	return &bchg{
		length:   2 + instSize,
		src:      src,
		register: register,
	}
}

func (b *bchg) Name() string {
	return fmt.Sprintf(
		"bchg    D%d, %s",
		b.register,
		b.src.string(),
	)
}

func (b *bchg) Length() int {
	return b.length
}
