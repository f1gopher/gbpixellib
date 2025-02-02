package instructions

import "fmt"

type moveToSR struct {
	size    instSize
	srcMode addressMode
	srcReg  byte
	data    uint32
}

func CreateMOVEToSR(value uint16, code []uint8, offset int) *moveToSR {
	size := wordSize
	destMode, destReg := effective(value, 0, size)
	return &moveToSR{
		size:    size,
		srcMode: destMode,
		srcReg:  destReg,
		data:    readData(code, offset, size),
	}
}

func (m *moveToSR) Name() string {
	return fmt.Sprintf("move    %s, SR", addressString(m.srcMode, m.srcReg, m.data, 0))
}

func (m *moveToSR) Length() int {
	return 2 + m.size.bytesUsed()
}
