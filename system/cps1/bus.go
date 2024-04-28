package cps1

import (
	"fmt"
	"log"

	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
)

type bus struct {
	log *log.Log

	rom    memory.Memory
	jamma  memory.Memory
	kick   memory.Memory
	gfxRam memory.Memory
	ram    memory.Memory
}

func CreateBus(log *log.Log) *bus {
	rom := make([]byte, 0x00400000)
	jamma := make([]byte, 0x00000018)
	kick := make([]byte, 0x00000002)
	gfxRam := make([]byte, 0x00930000)
	ram := make([]byte, 0x00010000)

	return &bus{
		log:    log,
		rom:    memory.CreateReadOnlyMemory("rom", rom, 0x00000000),
		jamma:  memory.CreateMemory("jamma", jamma, 0x00800000),
		kick:   memory.CreateMemory("kick", kick, 0x00800176),
		gfxRam: memory.CreateMemory("gfx", gfxRam, 0x00900000),
		ram:    memory.CreateMemory("ram", ram, 0x00FF0000),
	}
}

func (b *bus) Reset() {
	b.jamma.Reset()
	b.kick.Reset()
	// TODO - what else?
}

func (b *bus) ReadBit(address uint16, bit uint8) bool {
	return b.target(address).ReadBit(address, bit)
}

func (b *bus) ReadByte(address uint16) byte {
	return b.target(address).ReadByte(address)
}

func (b *bus) ReadShort(address uint16) uint16 {
	return b.target(address).ReadShort(address)
}

func (b *bus) WriteBit(address uint16, bit uint8, value bool) {
	b.target(address).WriteBit(address, bit, value)
}

func (b *bus) WriteByte(address uint16, value byte) {
	b.target(address).WriteByte(address, value)
}

func (b *bus) WriteShort(address uint16, value uint16) {
	b.target(address).WriteShort(address, value)
}

func (b *bus) LoadData(offset uint16, data []uint8) error {
	dest := b.target(offset)
	return dest.LoadData(offset, data)
}

func (b *bus) target(address uint16) memory.Memory {
	switch address {
	case address <= 0x003FFFFF:
		return b.rom
	case address >= 0x008000000 && address <= 0x00800037:
		return b.jamma
	case address >= 0x00800176 && address <= 0x00800177:
		return b.kick
	case address >= 0x00800100 && address <= 0x0080013F:
		panic("CPS-A not implemented")
	case address >= 0x00800140 && address <= 0x0080017F:
		panic("CPS-B not implemented")
	case address >= 0x00900000 && address <= 0x0092FFFF:
		return b.gfxRam
	case address >= 0x00FF00000 && address <= 0x00FFFFFF:
		return b.ram
	default:
		panic(fmt.Sprintf("Invalid address 0x%8X", address))
	}
}
