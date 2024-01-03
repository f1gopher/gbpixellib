package memory

import "fmt"

type cartridgeMBC1 struct {
	romBanks map[uint8]*Memory
	ramBanks map[uint8]*Memory

	ramEnable         uint8
	romBankNumber     uint8
	ramBankNumber     uint8
	bankingModeSelect uint8
	ramStart          uint16
}

func createCartridgeMBC1(romSize uint32, ramSize uint32, data *[]byte) Cartridge {

	if romSize >= M_1MiB {
		panic("Not implemented version of MBC1")
	}

	ram := make([]byte, ramSize)
	return &cartridgeMBC1{
		romBanks:          splitDataIntoBanks(0x0000, 0x4000, M_16Kb, data, "ROM"),
		ramBanks:          splitDataIntoBanks(0xA000, 0xA000, M_8Kb, &ram, "RAM"),
		ramEnable:         0x00,
		romBankNumber:     0x01,
		ramBankNumber:     0x00,
		bankingModeSelect: 0x00,
		ramStart:          0xA000,
	}
}

func (c *cartridgeMBC1) Reset() {
	for k, _ := range c.ramBanks {
		c.ramBanks[k].Reset()
	}
	// Don't reset the memory because that is the rom game

	c.ramEnable = 0x00
	c.romBankNumber = 0x01
	c.ramBankNumber = 0x00
	c.bankingModeSelect = 0x00
}

func (c *cartridgeMBC1) ReadBit(address uint16, bit uint8) bool {
	if address >= c.ramStart && !c.isRamEnabled() {
		return true
	}

	return c.memoryBank(address).ReadBit(address, bit)
}

func (c *cartridgeMBC1) ReadByte(address uint16) byte {
	if address >= c.ramStart && !c.isRamEnabled() {
		return 0xFF
	}

	return c.memoryBank(address).ReadByte(address)
}

func (c *cartridgeMBC1) ReadShort(address uint16) uint16 {
	if address >= c.ramStart && !c.isRamEnabled() {
		return 0xFFFF
	}

	return c.memoryBank(address).ReadShort(address)
}

func (c *cartridgeMBC1) WriteBit(address uint16, bit uint8, value bool) {
	if address <= 0x7FFF {
		panic("This shouldn't happen")
	}

	if address >= c.ramStart && !c.isRamEnabled() {
		return
	}

	c.memoryBank(address).WriteBit(address, bit, value)
}

func (c *cartridgeMBC1) WriteByte(address uint16, value byte) {
	if address >= c.ramStart && !c.isRamEnabled() {
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
	if address >= c.ramStart && address <= 0xBFFF && !c.isRamEnabled() {
		return
	}

	c.memoryBank(address).WriteByte(address, value)
}

func (c *cartridgeMBC1) WriteShort(address uint16, value uint16) {
	if address >= c.ramStart && !c.isRamEnabled() {
		return
	}

	if address <= 0x7FFF {
		panic("This shouldn't happen")
	}

	c.memoryBank(address).WriteShort(address, value)
}

func (c *cartridgeMBC1) DumpROMCode() (data []uint8, startAddress uint16) {

	// Combine the first and current banks
	code, startAddress := c.romBanks[0].DumpCode()
	currentCode, _ := c.romBanks[c.romBank()].DumpCode()
	code = append(code, currentCode...)

	return code, startAddress
}

func (c *cartridgeMBC1) DumpRAMCode() (data []uint8, startAddress uint16) {
	return c.ramBanks[c.ramBank()].DumpCode()
}

func (c *cartridgeMBC1) DumpROMBankCode(bank uint8) (data []uint8, startAddress uint16) {
	if int(bank) > len(c.romBanks)-1 {
		panic("Invalid bank number for cartridge")
	}
	return c.romBanks[bank].DumpCode()
}

func (c *cartridgeMBC1) DumpRAMBankCode(bank uint8) (data []uint8, startAddress uint16) {
	if int(bank) > len(c.ramBanks)-1 {
		panic("Invalid bank number for cartridge")
	}
	return c.ramBanks[bank].DumpCode()
}

func (c *cartridgeMBC1) CurrentROMBank() uint8 {
	return c.romBank()
}

func (c *cartridgeMBC1) CurrentRAMBank() uint8 {
	return c.ramBank()
}
func (c *cartridgeMBC1) memoryBank(address uint16) *Memory {
	if address >= c.ramStart && address <= 0xBFFF {
		bankNumber := c.ramBank()
		bank, exists := c.ramBanks[bankNumber]
		if !exists {
			panic(fmt.Sprintf("Cartridge RAM bank %d doesn't exist", bankNumber))
		}
		return bank
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

func (c *cartridgeMBC1) ramBank() uint8 {
	// TODO - is this right or should it be %?
	bank := (0b00000011 & c.ramBankNumber)
	return bank
}
