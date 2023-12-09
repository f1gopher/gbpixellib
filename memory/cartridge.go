package memory

type Cartridge struct {
	mem *Memory2
	ram *Memory2
}

func CreateCartridge(data *[]byte) *Cartridge {
	ram := make([]byte, 0x2000)
	return &Cartridge{
		mem: CreateMemory2("cartridge memory", data, 0),
		ram: CreateMemory2("cartridge ram", &ram, 0xA000),
	}
}

func (c *Cartridge) Reset() {
	c.ram.Reset()
	// Don't reset the memory because that is the rom game
}

func (c *Cartridge) ReadBit(address uint16, bit uint8) bool {
	return c.memoryBank(address).ReadBit(address, bit)
}

func (c *Cartridge) ReadByte(address uint16) byte {
	return c.memoryBank(address).ReadByte(address)
}

func (c *Cartridge) ReadShort(address uint16) uint16 {
	return c.memoryBank(address).ReadShort(address)
}

func (c *Cartridge) WriteBit(address uint16, bit uint8, value bool) {
	c.memoryBank(address).WriteBit(address, bit, value)
}

func (c *Cartridge) WriteByte(address uint16, value byte) {
	c.memoryBank(address).WriteByte(address, value)
}

func (c *Cartridge) WriteShort(address uint16, value uint16) {
	c.memoryBank(address).WriteShort(address, value)
}

func (c *Cartridge) DumpCode() []uint8 {
	return c.mem.DumpCode()
}

func (c *Cartridge) memoryBank(address uint16) *Memory2 {
	if address >= 0xA000 && address <= 0xBFFF {
		return c.ram
	}

	return c.mem
}
