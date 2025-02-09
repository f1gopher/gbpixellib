package instructions

import (
	"fmt"
)

type btst struct {
	size     instSize
	srcMode  addressMode
	srcReg   byte
	destMode addressMode
	destReg  byte

	srcData  uint32
	destData uint32
}

func CreateBTST(value uint16, code []uint8, offset int) *btst {
	size := zeroSize
	var srcData, destData uint32 = 0, 0
	srcMode := dataRegisterDirect
	srcReg := valueFrom3Bits(value, 9)
	destMode, destReg := effective(value, 0, wordSize)

	// bit number static
	if value&0b1111111111000000 == 0b0000100000000000 {
		srcMode = preIndexedImediate

		if destMode == addressRegisterIndirectWithDisplacement {
			size = longSize
			data := readData(code, offset, size)
			srcData = data >> 16
			destData = data & 0b00000000000000001111111111111111
		} else {
			size = wordSize
			data := readData(code, offset, size)
			srcData = data
			destData = 0
		}
	}

	return &btst{
		size:     size,
		srcMode:  srcMode,
		srcReg:   srcReg,
		destMode: destMode,
		destReg:  destReg,
		srcData:  srcData,
		destData: destData,
	}
}

func (b *btst) Name() string {
	return fmt.Sprintf(
		"btst    %s, %s",
		addressString(b.srcMode, b.srcReg, b.srcData, 0),
		addressString(b.destMode, b.destReg, b.destData, 0))
}

func (b *btst) Length() int {
	return 2 + b.size.bytesUsed()
}
