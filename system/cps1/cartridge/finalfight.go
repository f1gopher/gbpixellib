package cartridge

import (
	"os"
	"path"
)

type finalfight struct {
	romFolder string
}

func (f *finalfight) init(romFolder string) {
	f.romFolder = romFolder
}

func (f *finalfight) isValid() error {
	// TODO - check all files exist
	return nil
}

func (f *finalfight) load(memory systemMemory) error {
	// Load CPU
	f.loadFile(memory, "ff_36.11f", 0x0000)
	f.loadFile(memory, "ff_42.11h", 0x2000)
	f.loadFile(memory, "ff_37.12f", 0x4000)
	f.loadFile(memory, "ffe_43.12h", 0x6000)
	// TODO - word swap??
	f.loadFile(memory, "ff_32m.8h", 0x8000)

	// Gfx
	f.loadFile(memory, "ff-5m.7a", 0x0000)
	f.loadFile(memory, "ff-7m.9a", 0x0000)
	f.loadFile(memory, "ff-1m.3a", 0x0000)
	f.loadFile(memory, "ff-3m.5a", 0x0000)

	// Audio
	f.loadFile(memory, "ff-09.12b", 0x0000)

	// OKI samples
	f.loadFile(memory, "ff-18.11c", 0x0000)
	f.loadFile(memory, "ff-19.12c", 0x0000)

	// A board
	f.loadFile(memory, "buf1", 0x0000)
	f.loadFile(memory, "ioa1", 0x0000)
	f.loadFile(memory, "prg1", 0x0000)
	f.loadFile(memory, "rom1", 0x0000)
	f.loadFile(memory, "sou1", 0x0000)

	// B board
	f.loadFile(memory, "s224b.1a", 0x0000)
	f.loadFile(memory, "iob1.11e", 0x0000)
}

func (f *finalfight) loadFile(memory systemMemory, file string, address uint16) error {
	a := path.Join(f.romFolder, file)
	data, err := os.ReadFile(a)
	if err != nil {
		panic(err)
	}

	return memory.Load(address, data)
}
