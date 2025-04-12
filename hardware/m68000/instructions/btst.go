package instructions

import (
	"fmt"
)

type btst struct {
	length int
	size   instSize
	src    addressModeInfo
	dest   addressModeInfo
}

func CreateBTST(value uint16, code []uint8, codeOffset int) *btst {
	size := zeroSize
	src := addressModeInfo{
		mode: dataRegisterDirect,
		reg:  valueFrom3Bits(value, 9),
	}
	dest, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)

	// bit number static
	if value&0b1111111111000000 == 0b0000100000000000 {
		src.mode = preIndexedImediate

		// if dest.mode == addressRegisterIndirectWithDisplacement {
		size = wordSize
		src.data = readData(code, codeOffset, byteSize)
		// src.data = data >> 16
		// destData = data & 0b00000000000000001111111111111111
		dest.displacement = signExtend16BitTo32Bits(uint16(readData(code, codeOffset+2, wordSize)))
		// } else {
		// 	size = wordSize
		// 	data := readData(code, codeOffset, size)
		// 	src.data = data
		// 	dest.data = 0
		// }
		//
		if dest.mode == invalidAddressMode ||
			dest.mode == addressRegisterDirect ||
			dest.mode == preIndexedImediate {
			return nil
		}
	} else {
		if dest.mode == invalidAddressMode ||
			dest.mode == addressRegisterDirect {
			return nil
		}

		// Only want the frist byte in this mode
		if dest.mode == preIndexedImediate {
			dest.data = dest.data & 0x00FF
		}
	}

	return &btst{
		length: 2 + size.bytesUsed() + instSize,
		size:   size,
		src:    src,
		dest:   dest,
	}
}

func (b *btst) Name() string {
	return fmt.Sprintf(
		"BTST %s,%s",
		b.src.string(),
		b.dest.string())
}

func (b *btst) Length() int {
	return b.length
}
