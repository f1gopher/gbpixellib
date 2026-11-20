package m68000

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
)

type exg struct {
	opmode byte
	regX   byte
	regY   byte
}

const esgOpmode_DataRegs = 0b00001000
const esgOpmode_AddressRegs = 0b00001001
const esgOpmode_DataAddressRegs = 0b00010001

func CreateEXG(value uint16) *exg {
	opmode := byte((value & 0b0000000011111000) >> 3)
	regX := valueFrom3Bits(value, 9)
	regY := valueFrom3Bits(value, 0)

	return &exg{
		opmode: opmode,
		regX:   regX,
		regY:   regY,
	}
}

func (e *exg) Name() string {
	srcPrefix := "D"
	destPrefix := "D"
	if e.opmode == esgOpmode_AddressRegs {
		srcPrefix = "A"
		destPrefix = "A"
	} else if e.opmode == esgOpmode_DataAddressRegs {
		destPrefix = "A"
	}

	return fmt.Sprintf(
		"exg     %s%d, %s%d",
		srcPrefix,
		e.regX,
		destPrefix,
		e.regY)
}

func (e *exg) Length() int {
	return 2
}

func (e *exg) Step(registers *Registers, memory memory.Memory[uint32]) bool {
	switch e.opmode {
	case esgOpmode_AddressRegs:
		rx := registers.GetAReg32(e.regX)
		ry := registers.GetAReg32(e.regY)
		registers.SetAReg32(e.regX, ry)
		registers.SetAReg32(e.regY, rx)
	case esgOpmode_DataRegs:
		rx := registers.GetDReg32(e.regX)
		ry := registers.GetDReg32(e.regY)
		registers.SetDReg32(e.regX, ry)
		registers.SetDReg32(e.regY, rx)
	case esgOpmode_DataAddressRegs:
		rx := registers.GetDReg32(e.regX)
		ry := registers.GetAReg32(e.regY)
		registers.SetDReg32(e.regX, ry)
		registers.SetAReg32(e.regY, rx)
	default:
		panic(fmt.Errorf("Unhandled ESG opmode %08b", e.opmode))
	}

	return true
}
