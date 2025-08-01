package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type clr struct {
	opcode uint16
	cycle  uint8

	length int
	size   instSize
	dest   addressModeInfo
}

func CreateCLR(value uint16, code []uint8, codeOffset int) *clr {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	return &clr{
		length: 2 + instSize,
		size:   size,
		dest:   dest,
	}
}

func (c *clr) Name() string {
	return fmt.Sprintf(
		"clr.%s   %s",
		c.size.string(),
		c.dest.string())
}

func (c *clr) Length() int {
	return c.length
}

func (c *clr) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	switch c.cycle {
	case 0:
		size := twoBitSize(c.opcode, 6)
		dest, _, instSize := readEffectiveAddressMode(c.opcode, 0, size, memory, codeOffset, codeOffset)
		c.length = 2 + instSize
		c.size = size
		c.dest = dest

		c.cycle++
	case 1, 2:
		c.cycle++
	case 3:
		c.dest.setValueWithSize(registers, memory, 0, c.size)
		c.cycle++
		return true
	}

	return false
}
