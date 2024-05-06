package memoryBus

import (
	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type Bus struct {
	log *log.Log

	_68000 memory.Memory[uint32]
	gfxRam memory.Memory[uint32]
}

func CreateBus(log *log.Log) *Bus {
	ram := make([]uint8, 0x200000)

	return &Bus{
		log:    log,
		_68000: create68000(log),
		gfxRam: memory.CreateMemory[uint32]("gfx ram", &ram, 0x000000),
	}
}

func (b *Bus) Reset() {
	b._68000.Reset()
	b.gfxRam.Reset()
}

func (b *Bus) Memory68000() memory.Memory[uint32] {
	return b._68000
}

func (b *Bus) MemoryGfx() memory.Memory[uint32] {
	return b.gfxRam
}
