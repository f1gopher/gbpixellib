package memoryBus

import (
	"errors"
	"fmt"

	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type Bus struct {
	log *log.Log

	_68000 memory.Memory[uint32]
	gfxRam memory.Memory[uint32]
}

func CreateBus(log *log.Log) (bus *Bus, jammaMemory memory.Memory[uint32], kickMemory memory.Memory[uint32]) {
	_68000, gfx, jamma, kick := create68000(log)

	return &Bus{
		log:    log,
		_68000: _68000,
		gfxRam: gfx,
	}, jamma, kick
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

func (b *Bus) DumpMemory(address uint32) (data []uint8, err error) {
	err = nil
	switch {
	case address <= b._68000.LastAddress():
		data, err = b._68000.DumpRegionsCode(address)
	default:
		err = errors.New(fmt.Sprintf("Invalid address for memory dump: 0x%04X", address))
	}
	return
}
