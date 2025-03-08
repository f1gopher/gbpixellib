package instructions

import "fmt"

type bclr struct {
	length int
	size   instSize
	src    addressModeInfo
	dest   addressModeInfo

	data uint32
}

func CreateBCLR(value uint16, code []uint8, codeOffset int) *bclr {
	size := zeroSize
	var data uint32 = 0
	src := addressModeInfo{
		mode: dataRegisterDirect,
		reg:  valueFrom3Bits(value, 9),
	}

	// bit number static
	if value&0b1111111111000000 == 0b0000100010000000 {
		size = wordSize
		src.mode = preIndexedImediate
		src.data = readData(code, codeOffset, size)
	}

	dest, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

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
