package instructions

import "fmt"

type or struct {
	length  int
	size    instSize
	src     addressModeInfo
	destReg byte
}

func CreateOR(value uint16, code []uint8, codeOffset int) *or {
	opmode := valueFrom3Bits(value, 6)
	size := byteSize
	if opmode == 0b00000001 || opmode == 0b00000101 {
		size = wordSize
	} else if opmode == 0b00000010 || opmode == 0b00000110 {
		size = longSize
	}
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	destReg := valueFrom3Bits(value, 9)
	return &or{
		length:  2 + instSize,
		size:    size,
		src:     src,
		destReg: destReg,
	}
}

func (o *or) Name() string {
	return fmt.Sprintf(
		"or.%s    %s, D%d",
		o.size.string(),
		o.src.string(),
		o.destReg,
	)
}

func (o *or) Length() int {
	return o.length
}
