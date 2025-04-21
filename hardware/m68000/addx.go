package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type addx struct {
	size         instSize
	regX         byte
	regY         byte
	isAddressReg bool
}

func CreateADDX(value uint16, code []uint8, codeOffset int) *addx {
	size := twoBitSize(value, 6)

	if size == zeroSize {
		return nil
	}

	regX := valueFrom3Bits(value, 9)
	regY := valueFrom3Bits(value, 0)
	isAddressReg := value&0b0000000000001000 == 0b0000000000001000

	return &addx{
		size:         size,
		regX:         regX,
		regY:         regY,
		isAddressReg: isAddressReg,
	}
}

func (a *addx) Name() string {
	prefix := "D"
	if a.isAddressReg {
		prefix = "D"
	}

	return fmt.Sprintf(
		"addx.%s  %s%d, %s%d",
		a.size.string(),
		prefix,
		a.regY,
		prefix,
		a.regX)
}

func (a *addx) Length() int {
	return 2
}

func (a *addx) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
