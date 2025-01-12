package instructions

import "fmt"

type moveb struct {
	destMode addressMode
	destReg  byte
	srcMode  addressMode
	srcReg   byte
}

func CreateMOVEB(value uint16) *moveb {
	destMode, destReg := destination(value, 6)
	srcMode, srcReg := source(value, 0)

	return &moveb{
		destMode: destMode,
		destReg:  destReg,
		srcMode:  srcMode,
		srcReg:   srcReg,
	}
}

func (m *moveb) Name() string {
	return fmt.Sprintf("move.b  %s%d, %s%d", m.destMode.string(), m.destReg, m.srcMode.string(), m.srcReg)
}

func (m *moveb) Length() int {
	return 2
}
