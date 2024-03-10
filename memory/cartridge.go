package memory

import "fmt"

const M_1Kb = 1024
const M_8Kb = M_1Kb * 8
const M_16Kb = M_8Kb * 2
const M_32Kb = M_16Kb * 2
const M_1MiB = 1048576

type Cartridge interface {
	Reset()
	ReadBit(address uint16, bit uint8) bool
	ReadByte(address uint16) byte
	ReadShort(address uint16) uint16
	WriteBit(address uint16, bit uint8, value bool)
	WriteByte(address uint16, value byte)
	WriteShort(address uint16, value uint16)
	DumpROMCode() (data []uint8, startAddress uint16)
	DumpROMBankCode(bank uint8) (data []uint8, startAddress uint16)
	DumpRAMCode() (data []uint8, startAddress uint16)
	DumpRAMBankCode(bank uint8) (data []uint8, startAddress uint16)

	CurrentROMBank() uint8
	CurrentRAMBank() uint8
}

// Cartridge types and implementations
//
// 00h  ROM ONLY
// 01h  MBC1
// 02h  MBC1+RAM
// 03h  MBC1+RAM+BATTERY
// 05h  MBC2
// 06h  MBC2+BATTERY
// 08h  ROM+RAM
// 09h  ROM+RAM+BATTERY
// 0Bh  MMM01
// 0Ch  MMM01+RAM
// 0Dh  MMM01+RAM+BATTERY
// 0Fh  MBC3+TIMER+BATTERY
// 10h  MBC3+TIMER+RAM+BATTERY
// 11h  MBC3
// 12h  MBC3+RAM
// 13h  MBC3+RAM+BATTERY
// 19h  MBC5
// 1Ah  MBC5+RAM
// 1Bh  MBC5+RAM+BATTERY
// 1Ch  MBC5+RUMBLE
// 1Dh  MBC5+RUMBLE+RAM
// 1Eh  MBC5+RUMBLE+RAM+BATTERY
// 20h  MBC6
// 22h  MBC7+SENSOR+RUMBLE+RAM+BATTERY
// FCh  POCKET CAMERA
// FDh  BANDAI TAMA5
// FEh  HuC3
// FFh  HuC1+RAM+BATTERY

func IsCartridgeSupported(cartType uint8) bool {
	return cartType == 0x00 || // ROM Only
		cartType == 0x01 || // MBC1
		cartType == 0x02 || // MBC1+RAM
		cartType == 0x03 || // MBC1+RAM+BATTERY
		cartType == 0x05 || //MBC2
		cartType == 0x06 // MBC2+BATTERY
}

func CreateCartridge(cartType uint8, romSize uint32, ramSize uint32, data *[]byte) Cartridge {
	if len(*data) > int(romSize) {
		panic(fmt.Sprintf("Cartridge data size (%d) is bigger than the header specified rom size (%d)", len(*data), romSize))
	}

	if len(*data) != int(romSize) {
		panic("data doesn't match specified rom size")
	}

	switch cartType {
	// ROM Only
	case 0x00:
		return createCartridgeNoMBC(romSize, ramSize, data)
	// MBC1
	case 0x01, 0x02, 0x03:
		return createCartridgeMBC1(romSize, ramSize, data)
	// MBC2
	case 0x05, 0x06:
		return createCartridgeMBC2(romSize, ramSize, data)
	// MBC3
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		return createCartridgeMBC3(romSize, ramSize, data)
	default:
		panic(fmt.Sprintf("Unsupported cartridge type: 0x%02X", cartType))
	}
}
