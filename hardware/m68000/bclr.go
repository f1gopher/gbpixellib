package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type bclr struct {
	length int
	size   instSize
	src    addressModeInfo
	dest   addressModeInfo

	data uint32
}

func CreateBCLR(value uint16, code []uint8, codeOffset int) *bclr {
	offset := codeOffset
	size := zeroSize
	var data uint32 = 0
	src := addressModeInfo{
		mode: dataRegisterDirect,
		reg:  valueFrom3Bits(value, 9),
	}

	// bit number static
	if value&0b1111111111000000 == 0b0000100010000000 {
		size = byteSize
		src.mode = preIndexedImediate
		src.data = readData(code, codeOffset, size)
		offset += wordSize.bytesUsed()
	}

	dest, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, offset, codeOffset)

	if dest.mode == invalidAddressMode ||
		dest.mode == addressRegisterDirect ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement {
		return nil
	}

	return &bclr{
		length: 2 + instSize + size.bytesUsed(),
		size:   size,
		src:    src,
		dest:   dest,
		data:   data,
	}
}

func (b *bclr) Name() string {
	return fmt.Sprintf(
		"bclr    %s, %s",
		b.src.string(),
		b.dest.string())
}

func (b *bclr) Length() int {
	return b.length
}

func (b *bclr) Execute(registers *Registers, memory memory.Memory[uint32]) {
}
