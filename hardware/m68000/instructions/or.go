package instructions

import "fmt"

type or struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	destReg byte
}

func CreateOR(value uint16, code []uint8, offset int) *or {
	opmode := valueFrom3Bits(value, 6)
	size := byteSize
	if opmode == 0b00000001 || opmode == 0b00000101 {
		size = wordSize
	} else if opmode == 0b00000010 || opmode == 0b00000110 {
		size = longSize
	}
	srcMode, srcReg := effective(value, 0, size)
	destReg := valueFrom3Bits(value, 9)
	return &or{
		size:    size,
		srcMode: srcMode,
		srcReg:  srcReg,
		destReg: destReg,
	}
}

func (o *or) Name() string {
	return fmt.Sprintf(
		"or.%s    %s, D%d",
		o.size.string(),
		addressString(o.srcMode, o.srcReg, 0, 0),
		o.destReg,
	)
}

func (o *or) Length() int {
	return 2
}
