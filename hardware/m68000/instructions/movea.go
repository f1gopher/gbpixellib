package instructions

import "fmt"

type movea struct {
	src addressModeInfo

	opSize   instSize
	destReg  byte
	instSize int
}

func CreateMOVEA(value uint16, code []uint8, codeOffset int) *movea {
	size := twoBitSize(value, 12)
	// Use 6 as the offset even though it is 9 because we don't read the mode
	// so it will skip 6+3 to get to the destination
	destReg := valueFrom3Bits(value, 9)
	src, _, instSize := readSrcAddressMode(value, 0, size, code, codeOffset, codeOffset)

	return &movea{
		src:      src,
		opSize:   size,
		destReg:  destReg,
		instSize: 2 + instSize,
	}
}

func (m *movea) Name() string {
	return fmt.Sprintf("movea.%s %s, A%d", m.opSize.string(), m.src.string(), m.destReg)
}

func (m *movea) Length() int {
	return m.instSize
}
