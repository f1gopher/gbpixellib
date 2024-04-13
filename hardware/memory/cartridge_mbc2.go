package memory

import "fmt"

type cartridgeMBC2 struct {
	romBanks    map[uint8]*Memory
	ramBank     *Memory
	romBankMask uint8

	ramEnable         uint8
	romBankNumber     uint8
	bankingModeSelect uint8
	ramStart          uint16
}

func createCartridgeMBC2(romSize uint32, ramSize uint32, data *[]byte) Cartridge {
	ram := make([]byte, 512)
	return &cartridgeMBC2{
		romBanks:          splitDataIntoBanks(0x0000, 0x4000, M_16Kb, data, "ROM"),
		ramBank:           CreateReadOnlyMemory("RAM", &ram, 0xA000),
		romBankMask:       0b00000111,
		ramEnable:         0x00,
		romBankNumber:     0x01,
		bankingModeSelect: 0x00,
		ramStart:          0xA000,
	}
}

func (c *cartridgeMBC2) Reset() {
	c.ramBank.Reset()
	// Don't reset the memory because that is the rom game

	c.ramEnable = 0x00
	c.romBankNumber = 0x01
	c.bankingModeSelect = 0x00
}

func (c *cartridgeMBC2) ReadBit(address uint16, bit uint8) bool {
	if address >= c.ramStart && !c.isRamEnabled() {
		return true
	}

	return c.memoryBank(address).ReadBit(c.coerceAddressForShadow(address), bit)
}

func (c *cartridgeMBC2) ReadByte(address uint16) byte {
	if address >= c.ramStart && !c.isRamEnabled() {
		return 0xFF
	}

	return c.memoryBank(address).ReadByte(c.coerceAddressForShadow(address))
}

func (c *cartridgeMBC2) ReadShort(address uint16) uint16 {
	if address >= c.ramStart && !c.isRamEnabled() {
		return 0xFFFF
	}

	return c.memoryBank(address).ReadShort(c.coerceAddressForShadow(address))
}

func (c *cartridgeMBC2) WriteBit(address uint16, bit uint8, value bool) {
	if address <= 0x7FFF {
		panic("This shouldn't happen")
	}

	if address >= c.ramStart && !c.isRamEnabled() {
		return
	}

	c.memoryBank(address).WriteBit(c.coerceAddressForShadow(address), bit, value)
}

func (c *cartridgeMBC2) WriteByte(address uint16, value byte) {
	if address >= c.ramStart && !c.isRamEnabled() {
		return
	}

	if address <= 0x3FFF {
		// if address >= 0x0000 && address <= 0x1FFF {
		// 	c.ramEnable = value
		// } else if address >= 0x2000 && address <= 0x3FFF {
		// 	// ignore bits > 0x1F
		// 	value = 0b00011111 & value
		//
		// 	c.romBankNumber = value
		// } else if address >= 0x4000 && address <= 0x5FFF {
		// } else if address >= 0x6000 && address <= 0x7FFF {
		// 	c.bankingModeSelect = value
		// }
		if address&0b10000000 == 0b10000000 {
			c.romBankNumber = value
		} else {
			c.ramEnable = value
		}

		return
	}

	// If writing to RAM but RAM is disabled do nothing
	if address >= c.ramStart && address <= 0xBFFF && !c.isRamEnabled() {
		return
	}

	c.memoryBank(address).WriteByte(c.coerceAddressForShadow(address), value)
}

func (c *cartridgeMBC2) WriteShort(address uint16, value uint16) {
	if address >= c.ramStart && !c.isRamEnabled() {
		return
	}

	if address <= 0x7FFF {
		panic("This shouldn't happen")
	}

	c.memoryBank(address).WriteShort(c.coerceAddressForShadow(address), value)
}

func (c *cartridgeMBC2) DumpROMCode() (data []uint8, startAddress uint16) {

	// Combine the first and current banks
	code, startAddress := c.romBanks[0].DumpCode()
	currentCode, _ := c.romBanks[c.romBank()].DumpCode()
	code = append(code, currentCode...)

	return code, startAddress
}

func (c *cartridgeMBC2) DumpRAMCode() (data []uint8, startAddress uint16) {
	return c.ramBank.DumpCode()
}

func (c *cartridgeMBC2) DumpROMBankCode(bank uint8) (data []uint8, startAddress uint16) {
	if int(bank) > len(c.romBanks)-1 {
		panic("Invalid bank number for cartridge")
	}
	return c.romBanks[bank].DumpCode()
}

func (c *cartridgeMBC2) DumpRAMBankCode(bank uint8) (data []uint8, startAddress uint16) {
	return c.ramBank.DumpCode()
}

func (c *cartridgeMBC2) CurrentROMBank() uint8 {
	return c.romBank()
}

func (c *cartridgeMBC2) CurrentRAMBank() uint8 {
	return 0
}

func (c *cartridgeMBC2) memoryBank(address uint16) *Memory {
	if address >= c.ramStart && address <= 0xBFFF {
		// bankNumber := c.ramBank()
		// bank, exists := c.ramBanks[bankNumber]
		// if !exists {
		// 	panic(fmt.Sprintf("Cartridge RAM bank %d doesn't exist", bankNumber))
		// }
		// return bank
		return c.ramBank
	}

	// If the address is in the first bank that isn't switchable
	if address < 0x4000 {
		return c.romBanks[0]
	}

	bankNumber := c.romBank()
	bank, exists := c.romBanks[bankNumber]
	if !exists {
		panic(fmt.Sprintf("Cartridge ROM bank %d doesn't exist", bankNumber))
	}
	return bank
}

func (c *cartridgeMBC2) isRamEnabled() bool {
	// Lower 4 bits must be A
	return (0b00001111 & c.ramEnable) == 0x0A
}

func (c *cartridgeMBC2) romBank() uint8 {
	bank := (c.romBankMask & c.romBankNumber)

	if bank == 0 {
		bank = 1
	}

	return bank
}

func (c *cartridgeMBC2) coerceAddressForShadow(address uint16) uint16 {
	// Handle shadow RAM
	if address >= 0xA200 && address <= 0xBFFF {
		return address - 0x0200
	}

	return address
}

// func (c *cartridgeMBC2) ramBank() uint8 {
// 	// TODO - is this right or should it be %?
// 	bank := (0b00000011 & c.ramBankNumber)
// 	return bank
// }
