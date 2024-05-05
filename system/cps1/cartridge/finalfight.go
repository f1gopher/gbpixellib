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
	f.loadFilesInterlaced(memory, "ff_36.11f", "ff_42.11h", 0x00000000)
	f.loadFilesInterlaced(memory, "ff_37.12f", "ffe_43.12h", 0x00004000)
	// TODO - word swap??
	f.loadFile(memory, "ff_32m.8h", 0x00008000)

	// Gfx
	f.loadFile(memory, "ff-5m.7a", 0x00000000)
	f.loadFile(memory, "ff-7m.9a", 0x00000000)
	f.loadFile(memory, "ff-1m.3a", 0x00000000)
	f.loadFile(memory, "ff-3m.5a", 0x00000000)

	// Audio
	f.loadFile(memory, "ff-09.12b", 0x00000000)

	// OKI samples
	f.loadFile(memory, "ff-18.11c", 0x00000000)
	f.loadFile(memory, "ff-19.12c", 0x00000000)

	// A board
	f.loadFile(memory, "buf1", 0x00000000)
	f.loadFile(memory, "ioa1", 0x00000000)
	f.loadFile(memory, "prg1", 0x00000000)
	f.loadFile(memory, "rom1", 0x00000000)
	f.loadFile(memory, "sou1", 0x00000000)

	// B board
	f.loadFile(memory, "s224b.1a", 0x00000000)
	f.loadFile(memory, "iob1.11e", 0x00000000)

	return nil
}

func (f *finalfight) loadFile(memory systemMemory, file string, address uint32) error {
	a := path.Join(f.romFolder, file)
	data, err := os.ReadFile(a)
	if err != nil {
		panic(err)
	}

	return memory.LoadData(address, data)
}

func (f *finalfight) loadFilesInterlaced(memory systemMemory, file1 string, file2 string, address uint32) error {
	data1, err := os.ReadFile(path.Join(f.romFolder, file1))
	if err != nil {
		panic(err)
	}
	data2, err := os.ReadFile(path.Join(f.romFolder, file2))
	if err != nil {
		panic(err)
	}

	result := make([]uint8, len(data1)+len(data2))

	// Take alternating bytes from each file to assemble the complete buffer
	for x := 0; x < len(data1); x++ {
		result[x*2] = data1[x]
		result[(x*2)+1] = data2[x]
	}

	return memory.LoadData(address, result)
}
