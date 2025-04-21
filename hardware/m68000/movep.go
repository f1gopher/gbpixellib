package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type movep struct {
	length       int
	size         instSize
	opmode       byte
	addressReg   byte
	dataReg      byte
	displacement uint16
}

func CreateMOVEP(value uint16, code []uint8, codeOffset int) *movep {
	opmode := valueFrom3Bits(value, 6)
	addressReg := valueFrom3Bits(value, 0)
	dataReg := valueFrom3Bits(value, 9)
	size := wordSize

	if opmode == 0b00000100 || opmode == 0b00000110 {
		size = wordSize
	} else if opmode == 0b000000101 || opmode == 0b00000111 {
		size = longSize
	} else {
		return nil
	}
	displacement := readData(code, codeOffset, wordSize)

	return &movep{
		length:       2 + wordSize.bytesUsed(),
		size:         size,
		opmode:       opmode,
		addressReg:   addressReg,
		dataReg:      dataReg,
		displacement: uint16(displacement),
	}
}

func (m *movep) Name() string {
	if m.opmode == 0b00000100 || m.opmode == 0b00000101 {
		return fmt.Sprintf(
			"movep.%s ($%x,A%d), D%d",
			m.size.string(),
			m.displacement,
			m.addressReg,
			m.dataReg)
	}

	return fmt.Sprintf(
		"movep.%s D%d, ($%x,A%d)",
		m.size.string(),
		m.dataReg,
		m.displacement,
		m.addressReg)
}

func (m *movep) Length() int {
	return m.length
}

func (m *movep) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	return true
}
