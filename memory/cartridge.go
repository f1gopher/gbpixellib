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

func CreateCartridge(cartType uint8, romSize uint32, ramSize uint32, data *[]byte) Cartridge {
	if len(*data) > int(romSize) {
		panic(fmt.Sprintf("Cartridge data size (%d) is bigger than the header specified rom size (%d)", len(*data), romSize))
	}

	if len(*data) != int(romSize) {
		panic("data doesn't match specified rom size")
	}

	switch cartType {
	case 0x00:
		return createCartridgeNoMBC(romSize, ramSize, data)
	case 0x01:
		return createCartridgeMBC1(romSize, ramSize, data)
	case 0x03:
		return createCartridgeMBC3(romSize, ramSize, data)
	default:
		panic(fmt.Sprintf("Unsupported cartridge type: 0x%02X", cartType))
	}
}
