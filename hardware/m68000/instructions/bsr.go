package instructions

import "fmt"

type bsr struct {
	size         instSize
	displacement int

	tmpStartOffset int
}

func CreateBSR(value uint16, code []uint8, offset int) *bsr {
	size := byteSize
	displacement := twosComplementToInt8Bit(uint8(value & 0b0000000011111111))
	if displacement == 0 {
		size = wordSize
		displacement = twosComplementToInt16Bit(uint16(readData(code, offset, wordSize)))
	}

	return &bsr{
		size:           size,
		displacement:   displacement,
		tmpStartOffset: offset,
	}
}

func (b *bsr) Name() string {
	return fmt.Sprintf(
		"bsr     $%x",
		b.tmpStartOffset+b.displacement)
}

func (b *bsr) Length() int {
	return 2 + b.size.bytesUsed()
}
