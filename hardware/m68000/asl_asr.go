package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type asl_asr struct {
	length        int
	direction     direction
	isMemoryShift bool

	// Register Shift
	size                 instSize
	reg                  byte
	countReg             byte
	isRegisterShiftCount bool

	// Memory Shift
	src addressModeInfo
}

func CreateASL_ASR(value uint16, code []uint8, codeOffset int) *asl_asr {
	isMemoryShift := value&0b1111111011000000 == 0b1110000011000000
	direction := right
	if value&0b0000000100000000 == 0b0000000100000000 {
		direction = left
	}

	if isMemoryShift {
		src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

		return &asl_asr{
			length:        2 + instSize,
			isMemoryShift: isMemoryShift,
			direction:     direction,
			src:           src,
		}
	}

	size := twoBitSize(value, 6)
	reg := valueFrom3Bits(value, 0)
	countReg := valueFrom3Bits(value, 9)
	isRegisterShiftCount := value&0b0000000000100000 == 0b0000000000100000
	if !isRegisterShiftCount && countReg == 0 {
		countReg = 8
	}

	return &asl_asr{
		length:               2,
		isMemoryShift:        isMemoryShift,
		direction:            direction,
		size:                 size,
		reg:                  reg,
		countReg:             countReg,
		isRegisterShiftCount: isRegisterShiftCount,
	}
}

func (a *asl_asr) Name() string {
	dir := "l"
	if a.direction == right {
		dir = "r"
	}

	if a.isMemoryShift {
		return fmt.Sprintf(
			"as%s.%s   %s",
			dir,
			a.size.string(),
			a.src.string())
	}

	srcPrefix := "#"
	if a.isRegisterShiftCount {
		srcPrefix = "D"
	}
	return fmt.Sprintf(
		"as%s.%s   %s%d, D%d",
		dir,
		a.size.string(),
		srcPrefix,
		a.countReg,
		a.reg)
}

func (a *asl_asr) Length() int {
	return a.length
}

func (a *asl_asr) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
