package memory

type videoRam struct {
	mem *Memory
}

func CreateVideoRam() *videoRam {
	data := make([]byte, 0x2000)
	return &videoRam{
		mem: CreateMemory("video ram", &data, 0x8000),
	}
}

func (v *videoRam) Reset() {
	v.mem.Reset()
}

func (v *videoRam) ReadBit(address uint16, bit uint8) bool {
	return v.mem.ReadBit(address, bit)
}

func (v *videoRam) ReadByte(address uint16) byte {
	return v.mem.ReadByte(address)
}

func (v *videoRam) ReadShort(address uint16) uint16 {
	return v.mem.ReadShort(address)
}

func (v *videoRam) WriteBit(address uint16, bit uint8, value bool) {
	v.mem.WriteBit(address, bit, value)
}

func (v *videoRam) WriteByte(address uint16, value byte) {
	v.mem.WriteByte(address, value)
}

func (v *videoRam) WriteShort(address uint16, value uint16) {
	v.mem.WriteShort(address, value)
}
