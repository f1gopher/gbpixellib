package memory

type inputOutput interface {
	ReadDirectional() uint8
	ReadStandard() uint8
}

type interupt interface {
	TriggerTimerOverflow()
}

const DividerRegister = 0xFF04

type ram struct {
	mem *Memory

	io       inputOutput
	interupt interupt
}

func CreateRam() *ram {
	data := make([]byte, ramSize)
	return &ram{
		mem: CreateMemory("ram", &data, ramOffset),
		io:  nil,
	}
}

func (r *ram) SetIO(io inputOutput, interupt interupt) {
	r.io = io
	r.interupt = interupt
}

func (r *ram) Reset() {
	r.mem.Reset()
	// For controller
	r.mem.WriteByte(0xFF00, 0x3F)

	//  TODO - remove LCDC - to match for comaprisons
	//m.buffer[0xFF40] = 0x91
}

func (r *ram) ReadBit(address uint16, bit uint8) bool {
	return r.mem.ReadBit(address, bit)
}

func (r *ram) ReadByte(address uint16) byte {
	// Uncomment for comparison runs
	//	if address == 0xFF00 {
	//		reg := r.mem.ReadByte(address)
	//		P14 := (reg >> 4) & 0x01
	//		P15 := (reg >> 5) & 0x01

	//		if P14 == 0 {
	//			return reg & ^r.io.ReadDirectional()
	//		}

	//		if P15 == 0 {
	//			return reg & ^r.io.ReadStandard()
	//		}

	//		return reg | 0x0F
	//	}

	return r.mem.ReadByte(address)
}

func (r *ram) ReadShort(address uint16) uint16 {
	return r.mem.ReadShort(address)
}

func (r *ram) WriteBit(address uint16, bit uint8, value bool) {
	r.mem.WriteBit(address, bit, value)
}

func (r *ram) WriteByte(address uint16, value byte) {
	// If the CPU tries to write to the address (regardless of value) we
	// set the value to zero
	if address == 0xFF44 {
		//m.buffer[address] = 0
		r.mem.WriteByte(address, value)
		return
	}

	// Controller
	if address == 0xFF00 {
		P14 := (value >> 4) & 0x01
		P15 := (value >> 5) & 0x01

		current := 0xFF & (value | 0b11001111)

		if P14 == 0 {
			current &= r.io.ReadDirectional()
		}

		if P15 == 0 {
			current &= r.io.ReadStandard()
		}

		// Uncomment for comparison runs
		//current := (r.mem.ReadByte(address) & 0xCF) | (value & 0x30)

		r.mem.WriteByte(address, current)
		return
	}

	// LCD status register
	if address == 0xFF41 {
		reg := r.mem.ReadByte(address)
		// Keep the existing values for the first 3 bits
		value = SetBit(value, 0, reg&0b0000001 == 0b00000001)
		value = SetBit(value, 1, reg&0b0000010 == 0b00000010)
		value = SetBit(value, 2, reg&0b0000100 == 0b00000100)
		r.mem.WriteByte(address, value)
		return
	}

	r.mem.WriteByte(address, value)

	// Echo RAM - might not need this
	//
	// 0xE000 - 0xFDFF: Echo RAM
	//This section of memory directly mirrors the working RAM section - meaning if you write into the first address of working RAM (0xC000), the same value will appear in the first spot of echo RAM (0xE000). Nintendo actively discouraged developers from using this area of memory and as such we can just pretend it doesn't exist.
	if address >= 0xC000 && address <= 0xDDFF {
		if address+0x2000 <= 0xFDFF {
			r.mem.WriteByte(address+0x2000, value)
		}
	}
}

func (r *ram) WriteDividerRegister(value uint8) {
	r.WriteByte(DividerRegister, value)
}

func (r *ram) DisplaySetScanline(value uint8) {
	// Only used by the display
	r.mem.WriteByte(0xFF44, value)
}

func (r *ram) DisplaySetStatus(value uint8) {
	// Only display code uses this path
	r.mem.WriteByte(0xFF41, value)
}

func (r *ram) WriteShort(address uint16, value uint16) {
	r.mem.WriteShort(address, value)
}
