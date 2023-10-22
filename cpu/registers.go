package cpu

import (
	"fmt"

	"github.com/f1gopher/gbpixellib/memory"
)

type Register int

const (
	AF Register = iota
	BC
	DE
	HL
	SP
	PC
	A
	F
	B
	C
	D
	E
	H
	L
)

func (r Register) String() string {
	return [...]string{"AF", "BC", "DE", "HL", "SP", "PC", "A", "F", "B", "C", "D", "E", "H", "L"}[r]
}

type registerFlags int

const (
	ZFlag registerFlags = iota
	NFlag
	HFlag
	CFlag
)

func (f registerFlags) String() string {
	return [...]string{"Z", "N", "H", "C"}[f]
}

const zFlagBit = 7
const nFlagBit = 6
const hFlagBit = 5
const cFlagBit = 4

type Registers struct {
	regAF uint16
	regBC uint16
	regDE uint16
	regHL uint16
	regSP uint16
	regPC uint16

	imeEnabled bool
}

func (r *Registers) reset() {
	r.regAF = 0x0000
	r.regBC = 0x0000
	r.regDE = 0x0000
	r.regHL = 0x0000
	r.regSP = 0x0000
	r.regPC = 0x0000
	r.imeEnabled = false
}

func (r *Registers) SetIME(enabled bool) {
	r.imeEnabled = enabled
}

func (r *Registers) GetIME() bool {
	return r.imeEnabled
}

func (r *Registers) Get8(source Register) uint8 {
	switch source {
	case A:
		return getHigh(&r.regAF)
	case F:
		return getLow(&r.regAF)
	case B:
		return getHigh(&r.regBC)
	case C:
		return getLow(&r.regBC)
	case D:
		return getHigh(&r.regDE)
	case E:
		return getLow(&r.regDE)
	case H:
		return getHigh(&r.regHL)
	case L:
		return getLow(&r.regHL)
	default:
		panic(fmt.Sprintf("Not a valid 8bit register for reading: %s", source.String()))
	}
}

func (r *Registers) Get16Msb(source Register) uint8 {
	switch source {
	case AF:
		return r.Get8(A)
	case BC:
		return r.Get8(B)
	case DE:
		return r.Get8(D)
	case HL:
		return r.Get8(H)
	case SP:
		return Msb(r.Get16(SP))
	case PC:
		return Msb(r.Get16(PC))
	default:
		panic(fmt.Sprintf("Not a valid 16bit register for reading MSB: %s", source.String()))
	}
}

func (r *Registers) Get16Lsb(source Register) uint8 {
	switch source {
	case AF:
		return r.Get8(F)
	case BC:
		return r.Get8(C)
	case DE:
		return r.Get8(E)
	case HL:
		return r.Get8(L)
	case SP:
		return Lsb(r.Get16(SP))
	case PC:
		return Lsb(r.Get16(PC))
	default:
		panic(fmt.Sprintf("Not a valid 16bit register for reading LSB: %s", source.String()))
	}
}

func (r *Registers) Get16(source Register) uint16 {
	switch source {
	case AF:
		return r.regAF
	case BC:
		return r.regBC
	case DE:
		return r.regDE
	case HL:
		return r.regHL
	case SP:
		return r.regSP
	case PC:
		return r.regPC
	default:
		panic(fmt.Sprintf("Not a valid 16bit register for reading: %s", source.String()))
	}
}

func (r *Registers) Set8(target Register, value uint8) {
	switch target {
	case A:
		setHigh(&r.regAF, value)
	case F:
		setLow(&r.regAF, value)
	case B:
		setHigh(&r.regBC, value)
	case C:
		setLow(&r.regBC, value)
	case D:
		setHigh(&r.regDE, value)
	case E:
		setLow(&r.regDE, value)
	case H:
		setHigh(&r.regHL, value)
	case L:
		setLow(&r.regHL, value)
	default:
		panic(fmt.Sprintf("Not a valid 8bit register for writing: %s", target.String()))
	}
}

func (r *Registers) Set16FromTwoBytes(target Register, msb uint8, lsb uint8) {
	r.Set16(target, combineBytes(msb, lsb))
}

func (r *Registers) SetPC(value uint16) {
	r.Set16(PC, value)
}

func (r *Registers) Set16(target Register, value uint16) {
	switch target {
	case AF:
		r.regAF = value
	case BC:
		r.regBC = value
	case DE:
		r.regDE = value
	case HL:
		r.regHL = value
	case SP:
		r.regSP = value
	case PC:
		r.regPC = value
	default:
		panic(fmt.Sprintf("Not a valid 16bit register for writing: %s", target.String()))
	}
}

func (r *Registers) GetFlag(flag registerFlags) bool {
	switch flag {
	case ZFlag:
		return getRegBit(r.Get8(F), zFlagBit)
	case NFlag:
		return getRegBit(r.Get8(F), nFlagBit)
	case HFlag:
		return getRegBit(r.Get8(F), hFlagBit)
	case CFlag:
		return getRegBit(r.Get8(F), cFlagBit)
	default:
		panic("Unhandled flag for get")
	}
}

func (r *Registers) SetFlag(flag registerFlags, value bool) {
	switch flag {
	case ZFlag:
		r.SetRegBit(F, zFlagBit, value)
	case NFlag:
		r.SetRegBit(F, nFlagBit, value)
	case HFlag:
		r.SetRegBit(F, hFlagBit, value)
	case CFlag:
		r.SetRegBit(F, cFlagBit, value)
	default:
		panic("Unhandled flag for set")
	}
}

func (r *Registers) SetRegBit(reg Register, bit uint8, value bool) {
	result := memory.SetBit(r.Get8(reg), bit, value)
	r.Set8(reg, result)
}

func getHigh(value *uint16) uint8 {
	return uint8(*value >> 8)
}

func getLow(value *uint16) uint8 {
	return uint8(*value)
}

func setHigh(target *uint16, value uint8) {
	*target = *target &^ 0xFF00
	var final uint16 = uint16(value) << 8
	*target = *target | final
}

func setLow(target *uint16, value uint8) {
	*target = *target &^ 0x00FF
	*target = *target | uint16(value)
}

func getRegBit(value uint8, bit int) bool {
	return memory.GetBit(value, bit)
}
