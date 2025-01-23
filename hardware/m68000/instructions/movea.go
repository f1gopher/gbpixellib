package instructions

import "fmt"

type movea struct {
	size    instSize
	destReg byte
	srcMode addressMode
	srcReg  byte

	data uint32
}

func CreateMOVEA(value uint16, code []uint8, offset int) *movea {
	// Use 6 as the offset even though it is 9 because we don't read the mode
	// so it will skip 6+3 to get to the destination
	_, destReg := destination(value, 6)
	srcMode, srcReg := source(value, 0)
	size := twoBitSize(value, 12)

	return &movea{
		size:    size,
		data:    readData(code, offset, size),
		destReg: destReg,
		srcMode: srcMode,
		srcReg:  srcReg,
	}
}

func (m *movea) Name() string {
	return fmt.Sprintf("movea.%s #$%4x, A%d", m.size.string(), m.data, m.destReg)
}

func (m *movea) Length() int {
	return 2 + m.size.bytesUsed()
}
