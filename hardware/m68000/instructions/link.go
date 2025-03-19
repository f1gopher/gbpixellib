package instructions

import "fmt"

type link struct {
	reg          byte
	displacement uint16
}

func CreateLINK(value uint16, code []uint8, codeOffset int) *link {
	return &link{
		reg:          valueFrom3Bits(value, 0),
		displacement: uint16(readData(code, codeOffset, wordSize)),
	}
}

func (l *link) Name() string {
	return fmt.Sprintf(
		"link    A%d, #$%x",
		l.reg,
		l.displacement)
}

func (l *link) Length() int {
	return 4
}
