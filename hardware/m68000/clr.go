package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type clr struct {
	opcode uint16
	cycle  uint8

	length uint32
	size   instSize
	dest   addressModeInfo
}

func CreateCLR(value uint16, code []uint8, codeOffset int) *clr {
	size := twoBitSize(value, 6)
	dest, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	return &clr{
		length: uint32(2 + instSize),
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
	return int(c.length)
}

func (c *clr) StepClockCycle(registers *Registers, memory memory.Memory[uint32], pc uint32, prefetch prefetchQueue) (finished bool, newPC uint32) {
	switch c.cycle {
	case 0:
		size := twoBitSize(c.opcode, 6)
		dest, _, instSize := readEffectiveAddressMode2(c.opcode, 0, size, memory, pc, pc)
		c.length = 2 + instSize
		c.size = size
		c.dest = dest

		c.cycle++
		return false, pc
	case 1:
		c.cycle++
		return false, pc

	case 2:
		c.cycle++

		if c.dest.mode == dataRegisterDirect {
			c.dest.setValueWithSize(registers, memory, 0, c.size)
			return false, pc
		}
		// TODO - read value from prefetch
		c.dest.displacement = signExtend16BitTo32Bits(prefetch.ReadNextValue())
		return false, pc + 2

	case 3:
		c.cycle++
		if c.dest.mode == dataRegisterDirect {
			return true, pc
		}
		return false, pc
	case 4:
		c.cycle++
		// Read but discard the value
		regAddress := registers.GetAReg32(c.dest.reg)
		offset := displaceAddress(regAddress, c.dest.displacement)
		memory.ReadByte(offset)
		return false, pc
	case 5:
		c.cycle++
		return false, pc
	case 6:
		// Write takes 4 clock cycles/2 cpu cycles
		c.dest.setValueWithSize(registers, memory, 0, c.size)
		c.cycle++
		return false, pc
	case 7:
		c.cycle++
		return true, pc
	}

	return false, pc
}
