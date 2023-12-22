package memory

import "fmt"

type cartridgeMBC1 struct {
	romBanks map[uint8]*Memory
	ram      *Memory

	ramEnable         uint8
	romBankNumber     uint8
	ramBankNumber     uint8
	bankingModeSelect uint8
}

const bankSize = 0x4000

func splitDataIntoBanks(data *[]byte) map[uint8]*Memory {
	banks := make(map[uint8]*Memory)
	var currentBank uint8 = 0

	// Iterate the data and split into chunks of 0x4000 bytes per bank
	for x := 0; x < len(*data); x += bankSize {
		bank := make([]byte, bankSize)

		bank = (*data)[x : x+bankSize]

		// The first bank is always at 0 but other banks get swapped in and out
		// with the same off set of 0x4000
		var offset uint16 = 0
		if currentBank > 0 {
			offset = 0x4000
		}

		result := CreateReadOnlyMemory(fmt.Sprintf("cartridge ROM bank %d", currentBank), &bank, offset)
		banks[currentBank] = result
		currentBank++
	}

	return banks
}

func createCartridgeMBC1(romSize uint32, ramSize uint32, data *[]byte) Cartridge {

	if romSize >= M_1MiB {
		panic("Not implemented version of MBC1")
	}

	ram := make([]byte, ramSize)
	return &cartridgeMBC1{
		romBanks:          splitDataIntoBanks(data),
		ram:               CreateMemory("cartridge ram", &ram, 0xA000),
		ramEnable:         0x00,
		romBankNumber:     0x01,
		ramBankNumber:     0x00,
		bankingModeSelect: 0x00,
	}
}

func (c *cartridgeMBC1) Reset() {
	c.ram.Reset()
	// Don't reset the memory because that is the rom game

	c.ramEnable = 0x00
	c.romBankNumber = 0x01
	c.ramBankNumber = 0x00
	c.bankingModeSelect = 0x00
}

func (c *cartridgeMBC1) ReadBit(address uint16, bit uint8) bool {
	if address >= c.ram.addressOffset && !c.isRamEnabled() {
		return true
	}

	return c.memoryBank(address).ReadBit(address, bit)
}

func (c *cartridgeMBC1) ReadByte(address uint16) byte {
	if address >= c.ram.addressOffset && !c.isRamEnabled() {
		return 0xFF
	}

	return c.memoryBank(address).ReadByte(address)
}

func (c *cartridgeMBC1) ReadShort(address uint16) uint16 {
	if address >= c.ram.addressOffset && !c.isRamEnabled() {
		return 0xFFFF
	}

	return c.memoryBank(address).ReadShort(address)
}

func (c *cartridgeMBC1) WriteBit(address uint16, bit uint8, value bool) {
	if address <= 0x7FFF {
		panic("This shouldn't happen")
	}

	if address >= c.ram.addressOffset && !c.isRamEnabled() {
		return
	}

	c.memoryBank(address).WriteBit(address, bit, value)
}

func (c *cartridgeMBC1) WriteByte(address uint16, value byte) {
	if address >= c.ram.addressOffset && !c.isRamEnabled() {
		return
	}

	if address <= 0x7FFF {
		if address >= 0x0000 && address <= 0x1FFF {
			c.ramEnable = value
		} else if address >= 0x2000 && address <= 0x3FFF {
			// ignore bits > 0x1F
			value = 0b00011111 & value

			c.romBankNumber = value
		} else if address >= 0x4000 && address <= 0x5FFF {
			c.ramBankNumber = value
		} else if address >= 0x6000 && address <= 0x7FFF {
			c.bankingModeSelect = value
		}

		return
	}

	// If writing to RAM but RAM is disabled do nothing
	if address >= 0xA000 && address <= 0xBFFF && !c.isRamEnabled() {
		return
	}

	c.memoryBank(address).WriteByte(address, value)
}

func (c *cartridgeMBC1) WriteShort(address uint16, value uint16) {
	if address >= c.ram.addressOffset && !c.isRamEnabled() {
		return
	}

	if address <= 0x7FFF {
		panic("This shouldn't happen")
	}

	c.memoryBank(address).WriteShort(address, value)
}

func (c *cartridgeMBC1) DumpROMCode() []uint8 {

	// Combine the first and current banks
	code := c.romBanks[0].DumpCode()
	code = append(code, c.romBanks[c.romBank()].DumpCode()...)

	return code
}

func (c *cartridgeMBC1) DumpROMBankCode(bank uint8) []uint8 {
	if int(bank) > len(c.romBanks)-1 {
		panic("Invalid bank number for cartridge")
	}
	return c.romBanks[bank].DumpCode()
}

func (c *cartridgeMBC1) DumpRAMCode() []uint8 {
	return c.ram.DumpCode()
}

func (c *cartridgeMBC1) CurrentBank() uint8 {
	return c.romBank()
}

func (c *cartridgeMBC1) memoryBank(address uint16) *Memory {
	if address >= NoMBCRamOffset && address <= 0xBFFF {
		return c.ram
	}

	// If the address is in the first bank that isn't switchable
	if address < 0x4000 {
		return c.romBanks[0]
	}

	return c.romBanks[c.romBank()]
}

func (c *cartridgeMBC1) isRamEnabled() bool {
	// Lower 4 bits must be A
	return (0b00001111 & c.ramEnable) == 0x0A
}

func (c *cartridgeMBC1) romBank() uint8 {
	bank := (0b00011111 & c.romBankNumber)

	if bank == 0 {
		bank = 1
	}

	return bank
}
