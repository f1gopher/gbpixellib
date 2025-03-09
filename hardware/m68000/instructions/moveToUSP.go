package instructions

import "fmt"

type moveToUSP struct {
	reg            byte
	addressToStack bool
}

func CreateMOVEToUSP(value uint16) *moveToUSP {
	return &moveToUSP{
		reg:            valueFrom3Bits(value, 0),
		addressToStack: !(value&0b0000000000001000 == 0b0000000000001000),
	}
}

func (m *moveToUSP) Name() string {
	if m.addressToStack {
		return fmt.Sprintf("move    A%d, USP", m.reg)
	}
	return fmt.Sprintf("move    USP, A%d", m.reg)
}

func (m *moveToUSP) Length() int {
	return 2
}
