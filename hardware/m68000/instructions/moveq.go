package instructions

import "fmt"

type moveq struct {
	destReg byte

	data uint16
}

func CreateMOVEQ(value uint16) *moveq {
	destReg := uint8((value >> 9) & 0b0000000000000111)
	// TODO - "8 bits sign extended to a long"
	data := value & 0b0000000011111111

	return &moveq{
		destReg: destReg,
		data:    data,
	}
}

func (m *moveq) Name() string {
	return fmt.Sprintf("moveq   #$%x, D%d", m.data, m.destReg)
}

func (m *moveq) Length() int {
	return 2
}
