package instructions

import "fmt"

type moveusp struct {
	reg byte
}

func CreateMOVEUSP(value uint16) *moveusp {
	return &moveusp{
		reg: valueFrom3Bits(value, 0),
	}
}

func (m *moveusp) Name() string {
	return fmt.Sprintf("move    A%d, USP", m.reg)
}

func (m *moveusp) Length() int {
	return 2
}
