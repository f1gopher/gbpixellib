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
	mem *Memory2

	io       inputOutput
	interupt interupt
}

func CreateRam() *ram {
	data := make([]byte, ramSize)
	return &ram{
		mem: CreateMemory2("ram", &data, ramOffset),
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
	return r.mem.ReadByte(address)
}

func (r *ram) ReadShort(address uint16) uint16 {
	return r.mem.ReadShort(address)
}

func (r *ram) WriteBit(address uint16, bit uint8, value bool) {
	r.WriteBit(address, bit, value)
}

func (r *ram) WriteByte(address uint16, value byte) {
	// If the CPU tries to write to the address (regardless of value) we
	// set the value to zero
	if address == 0xFF44 {
		//m.buffer[address] = 0
		r.mem.WriteByte(address, value)
		return
	}

	// Any write resets it to 0
	if address == DividerRegister {
		r.mem.WriteByte(address, 0)
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

		r.mem.WriteByte(address, current)
		return
	}

	// Timer register
	if address == 0xFF05 {
		// When timer overflow from FF to 00 triogger interrupt
		if r.ReadByte(address) == 0xFF && value == 0x00 {
			r.interupt.TriggerTimerOverflow()
		}
	}

	//m.buffer[address] = value
	r.mem.WriteByte(address, value)

	// Echo RAM - might not need this
	//
	// 0xE000 - 0xFDFF: Echo RAM
	//This section of memory directly mirrors the working RAM section - meaning if you write into the first address of working RAM (0xC000), the same value will appear in the first spot of echo RAM (0xE000). Nintendo actively discouraged developers from using this area of memory and as such we can just pretend it doesn't exist.
	if address >= 0xC000 && address <= 0xDFFF {
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
	//m.buffer[0xFF44] = value
	r.mem.WriteByte(0xFF44, value)
}

func (r *ram) WriteShort(address uint16, value uint16) {
	r.mem.WriteShort(address, value)
}
