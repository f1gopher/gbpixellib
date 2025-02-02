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

	data uint32
}

func CreateBTST(value uint16, code []uint8, offset int) *btst {
	size := zeroSize
	var data uint32 = 0
	srcMode := dataRegisterDirect
	srcReg := valueFrom3Bits(value, 9)

	// bit number static
	if value&0b1111111111000000 == 0b0000100000000000 {
		size = wordSize
		srcMode = preIndexedImediate
		data = readData(code, offset, size)
	}

	destMode, destReg := effective(value, 0, wordSize)

	return &btst{
		size:     size,
		srcMode:  srcMode,
		srcReg:   srcReg,
		destMode: destMode,
		destReg:  destReg,
		data:     data,
	}
}

func (b *btst) Name() string {
	return fmt.Sprintf(
		"btst    %s, %s",
		addressString(b.srcMode, b.srcReg, b.data, 0),
		addressString(b.destMode, b.destReg, 0, 0))
}

func (b *btst) Length() int {
	return 2 + b.size.bytesUsed()
}
