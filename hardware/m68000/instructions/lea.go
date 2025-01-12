package instructions

import "fmt"

type lea struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	dstReg  byte

	data uint32
}

func CreateLEA(value uint16, code []uint8, offset int) *lea {
	size := wordSize
	srcMode, srcReg := effective(value, 0)
	dstReg := byte((value >> 9) & 0b0000000000000111)
	return &lea{
		size:    size,
		data:    readData(code, offset, size),
		srcMode: srcMode,
		srcReg:  srcReg,
		dstReg:  dstReg,
	}
}

func (l *lea) Name() string {
	return fmt.Sprintf("lea     $%4x.%s, %s%d", l.data, l.size.string(), l.srcMode.string(), l.dstReg)
}

func (l *lea) Length() int {
	return 2 + l.size.bytesUsed()
}
