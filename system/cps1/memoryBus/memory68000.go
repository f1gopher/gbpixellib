package memoryBus

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type memory68000 struct {
	log *log.Log

	rom    memory.Memory[uint32]
	jamma  memory.Memory[uint32]
	kick   memory.Memory[uint32]
	gfxRam memory.Memory[uint32]
	ram    memory.Memory[uint32]
}

func create68000(log *log.Log) (_680000 *memory68000, gfx memory.Memory[uint32], jammaMemory memory.Memory[uint32], kickMemory memory.Memory[uint32]) {
	rom := make([]byte, 0x00400000)
	jamma := make([]byte, 0x00000018)
	kick := make([]byte, 0x00000002)
	gfxRam := make([]byte, 0x00930000)
	ram := make([]byte, 0x00010000)
	gfxMemory := memory.CreateMemoryBigEndian[uint32]("gfx ram", &gfxRam, 0x900000)
	jammaMemory = memory.CreateMemoryBigEndian[uint32]("jamma", &jamma, JammaInputAddress)
	kickMemory = memory.CreateMemoryBigEndian[uint32]("kick", &kick, KickAddress)

	return &memory68000{
		log:    log,
		rom:    memory.CreateReadOnlyMemoryBigEndian[uint32]("rom", &rom, 0x000000),
		jamma:  jammaMemory,
		kick:   kickMemory,
		gfxRam: gfxMemory,
		ram:    memory.CreateMemoryBigEndian[uint32]("ram", &ram, 0xFF0000),
	}, gfxMemory, jammaMemory, kickMemory
}

func (m *memory68000) Jamma() memory.Memory[uint32] {
	return m.jamma
}

func (m *memory68000) Kick() memory.Memory[uint32] {
	return m.kick
}

func (m *memory68000) Reset() {
	m.jamma.Reset()
	m.kick.Reset()
	m.gfxRam.Reset()
	m.ram.Reset()
	// TODO - what else?
}

func (m *memory68000) Size() int {
	return int(m.ram.Offset()) + m.ram.Size()
}

func (m *memory68000) Offset() uint32 {
	return 0
}

func (m *memory68000) LastAddress() uint32 {
	return m.Offset() + uint32(m.Size()-1)
}

func (m *memory68000) ReadBit(address uint32, bit uint8) bool {
	return m.target(address).ReadBit(address, bit)
}

func (m *memory68000) ReadByte(address uint32) byte {
	return m.target(address).ReadByte(address)
}

func (m *memory68000) ReadShort(address uint32) uint16 {
	return m.target(address).ReadShort(address)
}

func (m *memory68000) WriteBit(address uint32, bit uint8, value bool) {
	m.target(address).WriteBit(address, bit, value)
}

func (m *memory68000) WriteByte(address uint32, value byte) {
	m.target(address).WriteByte(address, value)
}

func (m *memory68000) WriteShort(address uint32, value uint16) {
	m.target(address).WriteShort(address, value)
}

func (m *memory68000) DumpCode() (data []uint8, startAddress uint32) {
	dump := make([]uint8, m.Size())

	mem, startAddress := m.rom.DumpCode()
	copy(dump[startAddress:], mem)
	mem, startAddress = m.jamma.DumpCode()
	copy(dump[startAddress:], mem)
	mem, startAddress = m.kick.DumpCode()
	copy(dump[startAddress:], mem)
	mem, startAddress = m.gfxRam.DumpCode()
	copy(dump[startAddress:], mem)
	mem, startAddress = m.ram.DumpCode()
	copy(dump[startAddress:], mem)

	return dump, 0
}

func (m *memory68000) DumpRegionsCode(address uint32) (data []uint8, err error) {
	memory := m.target(address)
	if memory.Offset() != address {
		return nil, errors.New("Invalid dump address")
	}

	mem, _ := memory.DumpCode()
	return mem, nil
}

func (m *memory68000) LoadData(offset uint32, data []uint8) error {
	dest := m.target(offset)
	return dest.LoadData(offset, data)
}

func (m *memory68000) RawBuffer() []uint8 {
	panic("Not implemented")
}

func (m *memory68000) target(address uint32) memory.Memory[uint32] {
	switch {
	case address <= 0x003FFFFF:
		return m.rom
	case address >= 0x00800000 && address <= 0x00800037:
		return m.jamma
	case address >= 0x00800176 && address <= 0x00800177:
		return m.kick
	case address >= 0x00800100 && address <= 0x0080013F:
		panic("CPS-A not implemented")
	case address >= 0x00800140 && address <= 0x0080017F:
		panic("CPS-B.not implemented")
	case address >= 0x00900000 && address <= 0x0092FFFF:
		return m.gfxRam
	case address >= 0x00FF0000 && address <= 0x00FFFFFF:
		return m.ram
	default:
		panic(fmt.Sprintf("Invalid address 0x%8X", address))
	}
}
