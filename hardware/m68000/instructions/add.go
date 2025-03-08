package instructions

import "fmt"

type add struct {
	length  int
	opSize  instSize
	src     addressModeInfo
	mode    byte
	destReg byte
}

func CreateADD(value uint16, code []uint8, codeOffset int) *add {
	mode := valueFrom3Bits(value, 6)
	opSize := byteSize
	if mode == 0b001 || mode == 0b101 {
		opSize = wordSize
	} else if mode == 0b010 || mode == 0b110 {
		opSize = longSize
	}
	src, _, instSize := readEffectiveAddressMode(value, 0, opSize, code, codeOffset, codeOffset)
	destReg := valueFrom3Bits(value, 9)

	return &add{
		length:  2 + instSize,
		opSize:  opSize,
		src:     src,
		mode:    mode,
		destReg: destReg,
	}
}

func (a *add) Name() string {
	return fmt.Sprintf(
		"add.%s   %s, D%d",
		a.opSize.string(),
		a.src.string(),
		a.destReg,
	)
}

func (a *add) Length() int {
	return a.length
}
