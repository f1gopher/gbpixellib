package instructions

import "fmt"

type adda struct {
	length int
	size   instSize
	reg    byte
	src    addressModeInfo
}

func CreateADDA(value uint16, code []uint8, codeOffset int) *adda {
	size := wordSize
	opmode := valueFrom3Bits(value, 6)
	if opmode == 0b00000111 {
		size = longSize
	}
	reg := valueFrom3Bits(value, 9)
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	return &adda{
		length: 2 + instSize,
		size:   size,
		src:    src,
		reg:    reg,
	}
}

func (a *adda) Name() string {
	return fmt.Sprintf(
		"adda.%s  %s, A%d",
		a.size.string(),
		a.src.string(),
		a.reg)
}

func (a *adda) Length() int {
	return a.length
}
