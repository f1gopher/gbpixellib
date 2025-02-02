package instructions

import "fmt"

type lea struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	dstReg  byte

	data uint32

	// TODO - only needed for debug output, remove later and use PC value
	tmpStartOffset int
}

func CreateLEA(value uint16, code []uint8, offset int) *lea {
	size := wordSize
	srcMode, srcReg := effective(value, 0, size)
	dstReg := byte((value >> 9) & 0b0000000000000111)

	return &lea{
		size:           size,
		data:           readData(code, offset, size),
		srcMode:        srcMode,
		srcReg:         srcReg,
		dstReg:         dstReg,
		tmpStartOffset: offset,
	}
}

func (l *lea) Name() string {
	displacement := int(int16(uint16(l.data)))
	if l.srcMode == programCounterIndirectWithDisplacement {
		displacement = l.tmpStartOffset + displacement
	}

	return fmt.Sprintf("lea     %s, A%d",
		addressString(l.srcMode, l.srcReg, l.data, displacement),
		l.dstReg)
}

func (l *lea) Length() int {
	return 2 + l.size.bytesUsed()
}
