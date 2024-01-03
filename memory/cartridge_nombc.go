package memory

import "fmt"

type CartridgeNoMBC struct {
	rom *Memory
	ram *Memory
}

const NoMBCRamOffset = 0xA000

func createCartridgeNoMBC(romSize uint32, ramSize uint32, data *[]byte) *CartridgeNoMBC {
	ram := make([]byte, ramSize)
	return &CartridgeNoMBC{
		rom: CreateReadOnlyMemory("cartridge rom", data, 0),
		ram: CreateMemory("cartridge ram", &ram, NoMBCRamOffset),
	}
}

func (c *CartridgeNoMBC) Reset() {
	c.ram.Reset()
	// Don't reset the rom memory because that readonly
}

func (c *CartridgeNoMBC) ReadBit(address uint16, bit uint8) bool {
	return c.memoryBank(address).ReadBit(address, bit)
}

func (c *CartridgeNoMBC) ReadByte(address uint16) byte {
	return c.memoryBank(address).ReadByte(address)
}

func (c *CartridgeNoMBC) ReadShort(address uint16) uint16 {
	return c.memoryBank(address).ReadShort(address)
}

func (c *CartridgeNoMBC) WriteBit(address uint16, bit uint8, value bool) {
	c.memoryBank(address).WriteBit(address, bit, value)
}

func (c *CartridgeNoMBC) WriteByte(address uint16, value byte) {
	c.memoryBank(address).WriteByte(address, value)
}

func (c *CartridgeNoMBC) WriteShort(address uint16, value uint16) {
	c.memoryBank(address).WriteShort(address, value)
}

func (c *CartridgeNoMBC) DumpROMCode() (data []uint8, startAddress uint16) {
	return c.rom.DumpCode()
}

func (c *CartridgeNoMBC) DumpRAMCode() (data []uint8, startAddress uint16) {
	return c.ram.DumpCode()
}

func (c *CartridgeNoMBC) DumpROMBankCode(bank uint8) (data []uint8, startAddress uint16) {
	panic("Cartridge type does not have memory banks")
}

func (c *CartridgeNoMBC) DumpRAMBankCode(bank uint8) (data []uint8, startAddress uint16) {
	panic("Cartridge type does not have memory banks")
}

func (c *CartridgeNoMBC) CurrentROMBank() uint8 {
	return 0
}

func (c *CartridgeNoMBC) CurrentRAMBank() uint8 {
	return 0
}

func (c *CartridgeNoMBC) memoryBank(address uint16) *Memory {
	if address >= NoMBCRamOffset {
		if int(address) >= int(NoMBCRamOffset)+c.ram.size() {
			panic(fmt.Sprintf("Invalid address for NoMBC cartridge: 0x%2X", address))
		}
		return c.ram
	}

	return c.rom
}
