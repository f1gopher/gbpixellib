package instructions

import "fmt"

type bclr struct {
	size     instSize
	srcMode  addressMode
	srcReg   byte
	destMode addressMode
	destReg  byte

	data uint32
}

func CreateBCLR(value uint16, code []uint8, offset int) *bclr {
	size := zeroSize
	var data uint32 = 0
	srcMode := dataRegisterDirect
	srcReg := valueFrom3Bits(value, 9)

	// bit number static
	if value&0b1111111111000000 == 0b0000100010000000 {
		size = wordSize
		srcMode = preIndexedImediate
		data = readData(code, offset, size)
	}

	destMode, destReg := effective(value, 0, wordSize)

	return &bclr{
		size:     size,
		srcMode:  srcMode,
		srcReg:   srcReg,
		destMode: destMode,
		destReg:  destReg,
		data:     data,
	}
}

func (b *bclr) Name() string {
	return fmt.Sprintf(
		"bclr    %s, %s",
		addressString(b.srcMode, b.srcReg, b.data, 0),
		addressString(b.destMode, b.destReg, 0, 0))
}

func (b *bclr) Length() int {
	return 2 + b.size.bytesUsed()
}
