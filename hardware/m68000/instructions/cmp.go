package instructions

import "fmt"

type cmp struct {
	length  int
	src     addressModeInfo
	destReg byte
	size    instSize
}

func CreateCMP(value uint16, code []uint8, codeOffset int) *cmp {
	opmode := valueFrom3Bits(value, 6)
	size := byteSize
	if opmode == 0b00000001 {
		size = wordSize
	} else if opmode == 0b00000010 {
		size = longSize
	}

	destReg := valueFrom3Bits(value, 9)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	return &cmp{
		length:  2 + instSize,
		src:     src,
		size:    size,
		destReg: destReg,
	}
}

func (c *cmp) Name() string {
	return fmt.Sprintf(
		"cmp.%s   %s, D%d",
		c.size.string(),
		c.src.string(),
		c.destReg)
}

func (c *cmp) Length() int {
	return c.length
}
