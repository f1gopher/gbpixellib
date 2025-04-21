package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type eor struct {
	length  int
	size    instSize
	src     addressModeInfo
	mode    byte
	destReg byte
	data    uint32
}

func CreateEOR(value uint16, code []uint8, codeOffset int) *eor {
	mode := valueFrom3Bits(value, 6)
	size := byteSize
	if mode == 0b001 || mode == 0b101 {
		size = wordSize
	} else if mode == 0b010 || mode == 0b110 {
		size = longSize
	}
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)
	destReg := valueFrom3Bits(value, 9)

	return &eor{
		length:  2 + instSize,
		size:    size,
		src:     src,
		mode:    mode,
		destReg: destReg,
		data:    readData(code, codeOffset, size),
	}
}

func (e *eor) Name() string {
	return fmt.Sprintf(
		"eor.%s   D%d, %s",
		e.size.string(),
		e.destReg,
		e.src.string())
}

func (e *eor) Length() int {
	return e.length
}

func (e *eor) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
