package instructions

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type addq struct {
	length       int
	size         instSize
	dest         addressModeInfo
	value        byte
	data         uint32
	displacement int32
}

func CreateADDQ(value uint16, code []uint8, codeOffset int) *addq {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	imediateValue := (value >> 9) & 0b00000111
	// 0 means 8
	if imediateValue == 0 {
		imediateValue = 8
	}

	return &addq{
		length: 2 + instSize,
		size:   size,
		dest:   dest,
		value:  byte(imediateValue),
	}
}

func (a *addq) Name() string {
	return fmt.Sprintf(
		"addq.%s  #%d, %s",
		a.size.string(),
		a.value,
		a.dest.string())
}

func (a *addq) Length() int {
	// TODO - not sure thi sis right but with long size it only read an extra word
	return a.length
}

func (a *addq) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
