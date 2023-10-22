package debugger

import "github.com/f1gopher/gbpixellib/memory"

type debugMemory struct {
	memory *memory.Memory
}

func (d *debugMemory) reset() {
	d.memory.Reset()
}

func (d *debugMemory) ReadBit(address uint16, bit uint8) bool {
	return d.memory.ReadBit(address, bit)
}

func (d *debugMemory) ReadByte(address uint16) uint8 {
	return d.memory.ReadByte(address)
}

func (d *debugMemory) ReadShort(address uint16) uint16 {
	return d.ReadShort(address)
}

func (d *debugMemory) WriteByte(address uint16, value uint8) {
	d.memory.WriteByte(address, value)
}

func (d *debugMemory) WriteShort(address uint16, value uint16) {
	d.WriteShort(address, value)
}

func (d *debugMemory) Write(address uint16, data []uint8) error {
	return d.memory.Write(address, data)
}
