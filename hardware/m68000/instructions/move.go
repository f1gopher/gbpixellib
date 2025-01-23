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
	return fmt.Sprintf(
		"move.%s  %s, %s",
		m.size.string(),
		addressString(m.srcMode, m.srcReg, m.data),
		addressString(m.destMode, m.destReg, m.data))
}

func (m *move) Length() int {
	if m.srcMode == absoluteShort || m.srcMode == absoluteLong ||
		m.destMode == absoluteShort || m.destMode == absoluteLong {
		return 2 + m.size.bytesUsed()
	}

	return 2
}
