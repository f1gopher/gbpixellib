package instructions

import "fmt"

type bset struct {
	length int
	size   instSize
	src    addressModeInfo
	dest   addressModeInfo
}

func CreateBSET(value uint16, code []uint8, codeOffset int) *bset {
	size := zeroSize
	src := addressModeInfo{
		mode: dataRegisterDirect,
		reg:  valueFrom3Bits(value, 9),
	}
	dest, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	// bit number static
	if value&0b1111111111000000 == 0b0000100011000000 {
		src.mode = preIndexedImediate

		if dest.mode == addressRegisterIndirectWithDisplacement {
			size = wordSize
			src.data = readData(code, codeOffset, size)
			// src.data = data >> 16
			// destData = data & 0b00000000000000001111111111111111
			dest.displacement = signExtend16BitTo32Bits(uint16(readData(code, codeOffset+2, wordSize)))
		} else {
			size = wordSize
			data := readData(code, codeOffset, size)
			src.data = data
			dest.data = 0
		}
	}

	if dest.mode == invalidAddressMode ||
		dest.mode == addressRegisterDirect ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement {
		return nil
	}

	return &bset{
		length: 2 + size.bytesUsed() + instSize,
		size:   size,
		src:    src,
		dest:   dest,
	}
}

func (b *bset) Name() string {
	return fmt.Sprintf(
		"BSET %s,%s",
		b.src.string(),
		b.dest.string())
}

func (b *bset) Length() int {
	return b.length
}
