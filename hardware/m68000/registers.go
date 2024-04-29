package m68000

type register int

const (
	A0 register = iota
	A1
	A2
	A3
	A4
	A5
	A6
	A7
	D0
	D1
	D2
	D3
	D4
	D5
	D6
	D7
	PC
	SR
)

type registers struct {
	d0 uint32
	d1 uint32
	d2 uint32
	d3 uint32
	d4 uint32
	d5 uint32
	d6 uint32
	d7 uint32

	a0 uint32
	a1 uint32
	a2 uint32
	a3 uint32
	a4 uint32
	a5 uint32
	a6 uint32
	a7 uint32

	pc uint32
	sr uint32
}

func (r *registers) reset() {
	r.d0 = 0x00000000
	r.d1 = 0x00000000
	r.d2 = 0x00000000
	r.d3 = 0x00000000
	r.d4 = 0x00000000
	r.d5 = 0x00000000
	r.d6 = 0x00000000
	r.d7 = 0x00000000

	r.a0 = 0x00000000
	r.a1 = 0x00000000
	r.a2 = 0x00000000
	r.a3 = 0x00000000
	r.a4 = 0x00000000
	r.a5 = 0x00000000
	r.a6 = 0x00000000
	r.a7 = 0x00000000

	r.pc = 0x00000000
	r.sr = 0x00000000
}

func (r *registers) get8(target register) uint8 {
	switch target {
	case A0:
		return getRegister[uint8](&r.a0)
	case A1:
		return getRegister[uint8](&r.a1)
	case A2:
		return getRegister[uint8](&r.a2)
	case A3:
		return getRegister[uint8](&r.a3)
	case A4:
		return getRegister[uint8](&r.a4)
	case A5:
		return getRegister[uint8](&r.a5)
	case A6:
		return getRegister[uint8](&r.a6)
	case A7:
		return getRegister[uint8](&r.a7)
	case D0:
		return getRegister[uint8](&r.d0)
	case D1:
		return getRegister[uint8](&r.d1)
	case D2:
		return getRegister[uint8](&r.d2)
	case D3:
		return getRegister[uint8](&r.d3)
	case D4:
		return getRegister[uint8](&r.d4)
	case D5:
		return getRegister[uint8](&r.d5)
	case D6:
		return getRegister[uint8](&r.d6)
	case D7:
		return getRegister[uint8](&r.d7)
	case PC:
		return getRegister[uint8](&r.pc)
	case SR:
		return getRegister[uint8](&r.sr)
	default:
		panic("Unhandled register")
	}
}

func (r *registers) get16(target register) uint16 {
	switch target {
	case A0:
		return getRegister[uint16](&r.a0)
	case A1:
		return getRegister[uint16](&r.a1)
	case A2:
		return getRegister[uint16](&r.a2)
	case A3:
		return getRegister[uint16](&r.a3)
	case A4:
		return getRegister[uint16](&r.a4)
	case A5:
		return getRegister[uint16](&r.a5)
	case A6:
		return getRegister[uint16](&r.a6)
	case A7:
		return getRegister[uint16](&r.a7)
	case D0:
		return getRegister[uint16](&r.d0)
	case D1:
		return getRegister[uint16](&r.d1)
	case D2:
		return getRegister[uint16](&r.d2)
	case D3:
		return getRegister[uint16](&r.d3)
	case D4:
		return getRegister[uint16](&r.d4)
	case D5:
		return getRegister[uint16](&r.d5)
	case D6:
		return getRegister[uint16](&r.d6)
	case D7:
		return getRegister[uint16](&r.d7)
	case PC:
		return getRegister[uint16](&r.pc)
	case SR:
		return getRegister[uint16](&r.sr)
	default:
		panic("Unhandled register")
	}
}

func (r *registers) get32(target register) uint32 {
	switch target {
	case A0:
		return getRegister[uint32](&r.a0)
	case A1:
		return getRegister[uint32](&r.a1)
	case A2:
		return getRegister[uint32](&r.a2)
	case A3:
		return getRegister[uint32](&r.a3)
	case A4:
		return getRegister[uint32](&r.a4)
	case A5:
		return getRegister[uint32](&r.a5)
	case A6:
		return getRegister[uint32](&r.a6)
	case A7:
		return getRegister[uint32](&r.a7)
	case D0:
		return getRegister[uint32](&r.d0)
	case D1:
		return getRegister[uint32](&r.d1)
	case D2:
		return getRegister[uint32](&r.d2)
	case D3:
		return getRegister[uint32](&r.d3)
	case D4:
		return getRegister[uint32](&r.d4)
	case D5:
		return getRegister[uint32](&r.d5)
	case D6:
		return getRegister[uint32](&r.d6)
	case D7:
		return getRegister[uint32](&r.d7)
	case PC:
		return getRegister[uint32](&r.pc)
	case SR:
		return getRegister[uint32](&r.sr)
	default:
		panic("Unhandled register")
	}
}

func (r *registers) set8(target register, value uint8) {
	switch target {
	case A0, A1, A2, A3, A4, A5, A6, A7:
		panic("Unsupported operation for an addres register")
	case D0:
		setRegister[uint8](&r.d0, value)
	case D1:
		setRegister[uint8](&r.d1, value)
	case D2:
		setRegister[uint8](&r.d2, value)
	case D3:
		setRegister[uint8](&r.d3, value)
	case D4:
		setRegister[uint8](&r.d4, value)
	case D5:
		setRegister[uint8](&r.d5, value)
	case D6:
		setRegister[uint8](&r.d6, value)
	case D7:
		setRegister[uint8](&r.d7, value)
	case PC:
		setRegister[uint8](&r.pc, value)
	case SR:
		setRegister[uint8](&r.sr, value)
	}
}

func (r *registers) set16(target register, value uint16) {
	switch target {
	case A0:
		setRegisterSignExtend(&r.a0, value)
	case A1:
		setRegisterSignExtend(&r.a1, value)
	case A2:
		setRegisterSignExtend(&r.a2, value)
	case A3:
		setRegisterSignExtend(&r.a3, value)
	case A4:
		setRegisterSignExtend(&r.a4, value)
	case A5:
		setRegisterSignExtend(&r.a5, value)
	case A6:
		setRegisterSignExtend(&r.a6, value)
	case A7:
		setRegisterSignExtend(&r.a7, value)
	case D0:
		setRegister[uint16](&r.d0, value)
	case D1:
		setRegister[uint16](&r.d1, value)
	case D2:
		setRegister[uint16](&r.d2, value)
	case D3:
		setRegister[uint16](&r.d3, value)
	case D4:
		setRegister[uint16](&r.d4, value)
	case D5:
		setRegister[uint16](&r.d5, value)
	case D6:
		setRegister[uint16](&r.d6, value)
	case D7:
		setRegister[uint16](&r.d7, value)
	case PC:
		setRegister[uint16](&r.pc, value)
	case SR:
		setRegister[uint16](&r.sr, value)
	}
}

func (r *registers) set32(target register, value uint32) {
	switch target {
	case A0:
		setRegister[uint32](&r.a0, value)
	case A1:
		setRegister[uint32](&r.a1, value)
	case A2:
		setRegister[uint32](&r.a2, value)
	case A3:
		setRegister[uint32](&r.a3, value)
	case A4:
		setRegister[uint32](&r.a4, value)
	case A5:
		setRegister[uint32](&r.a5, value)
	case A6:
		setRegister[uint32](&r.a6, value)
	case A7:
		setRegister[uint32](&r.a7, value)
	case D0:
		setRegister[uint32](&r.d0, value)
	case D1:
		setRegister[uint32](&r.d1, value)
	case D2:
		setRegister[uint32](&r.d2, value)
	case D3:
		setRegister[uint32](&r.d3, value)
	case D4:
		setRegister[uint32](&r.d4, value)
	case D5:
		setRegister[uint32](&r.d5, value)
	case D6:
		setRegister[uint32](&r.d6, value)
	case D7:
		setRegister[uint32](&r.d7, value)
	case PC:
		setRegister[uint32](&r.pc, value)
	case SR:
		setRegister[uint32](&r.sr, value)
	}
}

func getRegister[T uint8 | uint16 | uint32](reg *uint32) T {
	switch any(*reg).(type) {
	case uint8:
		return T(*reg)
	case uint16:
		return T(*reg)
	case uint32:
		return T(*reg)
	default:
		panic("Unhandled datatype")
	}
}

func setRegister[T uint8 | uint16 | uint32](reg *uint32, value T) {
	switch any(*reg).(type) {
	case uint8:
		*reg &^= 0x000000FF
		*reg |= uint32(value)
	case uint16:
		*reg &^= 0x0000FFFF
		*reg |= uint32(value)
	case uint32:
		*reg = uint32(value)
	default:
		panic("Unhandled datatype")
	}
}

func setRegisterSignExtend(reg *uint32, value uint16) {
	// If the 15bit is 1 then fill the upper 2 bytes with 1
	if value&0b1000000000000000 == 0b1000000000000000 {
		*reg = 0xFFFF0000
		*reg = *reg | uint32(value)
		return
	}

	// Default to 0's padding
	*reg = uint32(value)
}
