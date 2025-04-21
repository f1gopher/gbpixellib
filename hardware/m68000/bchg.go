package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type bchg struct {
	length      int
	src         addressModeInfo
	register    byte
	data        byte
	isImmediate bool
}

func CreateBCHG(value uint16, code []uint8, codeOffset int) *bchg {
	var register byte = 0
	var data byte = 0
	offset := codeOffset
	isImmediate := value&0b1111111111000000 == 0b0000100001000000
	length := 2

	if isImmediate {
		data = uint8(readData(code, codeOffset, byteSize))
		offset += wordSize.bytesUsed()
		length += wordSize.bytesUsed()
	} else {
		register = valueFrom3Bits(value, 9)
	}

	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, offset, codeOffset)

	return &bchg{
		length:      length + instSize,
		src:         src,
		register:    register,
		data:        data,
		isImmediate: isImmediate,
	}
}

func (b *bchg) Name() string {
	if b.isImmediate {
		return fmt.Sprintf(
			"bchg    #$%x, %s",
			b.data,
			b.src.string(),
		)
	}

	return fmt.Sprintf(
		"bchg    D%d, %s",
		b.register,
		b.src.string(),
	)
}

func (b *bchg) Length() int {
	return b.length
}

func (b *bchg) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
