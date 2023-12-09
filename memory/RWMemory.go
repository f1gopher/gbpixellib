package memory

type RWMemory interface {
	ReadBit(address uint16, bit uint8) bool
	ReadByte(address uint16) byte
	ReadShort(address uint16) uint16
	WriteBit(address uint16, bit uint8, value bool)
	WriteByte(address uint16, value byte)
	WriteShort(address uint16, value uint16)
}
