package instructions

import "fmt"

type move struct {
	size     instSize
	destMode addressMode
	destReg  byte
	srcMode  addressMode
	srcReg   byte

	data uint32
}

func CreateMOVE(value uint16, code []uint8, offset int) *move {
	destMode, destReg := destination(value, 6)
	srcMode, srcReg := source(value, 0)
	size := moveSize(value, 12)

	return &move{
		size:     size,
		destMode: destMode,
		destReg:  destReg,
		srcMode:  srcMode,
		srcReg:   srcReg,
		data:     readData(code, offset, size),
	}
}

func (m *move) Name() string {
	if m.srcMode == absoluteShort || m.srcMode == absoluteLong {
		return fmt.Sprintf("move.%s  %s%04x, %s%d", m.size.string(), m.srcMode.string(), m.data, m.destMode.string(), m.destReg)
	}

	return fmt.Sprintf("move.%s  %s%d, %s%d", m.size.string(), m.srcMode.string(), m.srcReg, m.destMode.string(), m.destReg)
}

func (m *move) Length() int {
	if m.srcMode == absoluteShort || m.srcMode == absoluteLong ||
		m.destMode == absoluteShort || m.destMode == absoluteLong {
		return 2 + m.size.bytesUsed()
	}

	return 2
}
