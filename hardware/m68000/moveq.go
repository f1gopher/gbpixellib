package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type InstructionRegisters interface {
	getFlag(flag flag) bool
}

type moveq struct {
	opcode uint16
	cycle  uint8

	destReg byte

	data int16
}

func CreateMOVEQ(value uint16) *moveq {
	destReg := uint8((value >> 9) & 0b0000000000000111)
	// TODO - "8 bits sign extended to a long"
	data := int16(int8(value & 0b0000000011111111))

	return &moveq{
		destReg: destReg,
		data:    data,
	}
}

func (m *moveq) Name() string {
	// Some fiddling to match the display for negative numbers
	if m.data < 0 {
		return fmt.Sprintf("moveq   #-$%x, D%d", m.data*-1, m.destReg)
	}

	return fmt.Sprintf("moveq   #$%x, D%d", m.data, m.destReg)
}

func (m *moveq) Length() int {
	return 2
}

func (m *moveq) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	switch m.cycle {
	case 0:
		m.destReg = uint8((m.opcode >> 9) & 0b0000000000000111)
		// TODO - "8 bits sign extended to a long"
		m.data = int16(int8(m.opcode & 0b0000000011111111))
		m.cycle++
	case 1:
		m.cycle++
	case 2:
		m.cycle++
	case 3:
		registers.SetDReg32(m.destReg, uint32(m.data))
		m.cycle++
		return true
	}

	return false
}
