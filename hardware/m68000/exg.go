package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type exg struct {
	opmode byte
	regX   byte
	regY   byte
}

func CreateEXG(value uint16, code []uint8, codeOffset int) *exg {
	opmode := byte((value & 0b0000000011111000) >> 3)
	regX := valueFrom3Bits(value, 9)
	regY := valueFrom3Bits(value, 0)

	return &exg{
		opmode: opmode,
		regX:   regX,
		regY:   regY,
	}
}

func (e *exg) Name() string {
	srcPrefix := "D"
	destPrefix := "D"
	if e.opmode == 0b00001001 {
		srcPrefix = "A"
		destPrefix = "A"
	} else if e.opmode == 0b00010001 {
		destPrefix = "A"
	}

	return fmt.Sprintf(
		"exg     %s%d, %s%d",
		srcPrefix,
		e.regX,
		destPrefix,
		e.regY)
}

func (e *exg) Length() int {
	return 2
}

func (e *exg) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
