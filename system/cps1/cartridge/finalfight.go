package cartridge

import (
	"os"
	"path"

	"github.com/f1gopher/gbpixellib/hardware/memory"
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

func (f *finalfight) load(_68000 memory.Memory[uint32], gfx memory.Memory[uint32]) error {
	// Load CPU
	f.loadFilesInterlaced(_68000, "ff_36.11f", "ff_42.11h", 0x00000000)
	f.loadFilesInterlaced(_68000, "ff_37.12f", "ffe_43.12h", 0x00040000)
	f.loadFileSwapUint16(_68000, "ff-32m.8h", 0x00080000)

	// Gfx
	f.loadFilesInterlacedMad(gfx, "ff-5m.7a", "ff-7m.9a", "ff-1m.3a", "ff-3m.5a", 0x00000000)

	// Audio
	// f.loadFile(memory, "ff_09.12b", 0x00000000)

	// OKI samples
	// f.loadFile(memory, "ff_18.11c", 0x00000000)
	// f.loadFile(memory, "ff_19.12c", 0x00020000)

	// A board
	// TODO - this isn't right - same as B board
	// f.loadFile(memory, "buf1", 0x00000000)
	// f.loadFile(memory, "ioa1", 0x00000000)
	// f.loadFile(memory, "prg1", 0x00000000)
	// f.loadFile(memory, "rom1", 0x00000000)
	// f.loadFile(memory, "sou1", 0x00000000)

	// B board
	// TODO - load one file over the top of the other and last file wins???!
	// has an extra 0A at the end for some reason
	// f.loadFile(memory, "s224b.1a", 0x00000000)
	// f.loadFile(memory, "iob1.11e", 0x00000000)

	return nil
}

func (f *finalfight) loadFile(memory memory.Memory[uint32], file string, address uint32) error {
	a := path.Join(f.romFolder, file)
	data, err := os.ReadFile(a)
	if err != nil {
		panic(err)
	}

	return memory.LoadData(address, data)
}

func (f *finalfight) loadFilesInterlaced(memory memory.Memory[uint32], file1 string, file2 string, address uint32) error {
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

func (f *finalfight) loadFilesInterlacedMad(memory memory.Memory[uint32], file1 string, file2 string, file3 string, file4 string, address uint32) error {
	data1, err := os.ReadFile(path.Join(f.romFolder, file1))
	if err != nil {
		panic(err)
	}
	data2, err := os.ReadFile(path.Join(f.romFolder, file2))
	if err != nil {
		panic(err)
	}
	data3, err := os.ReadFile(path.Join(f.romFolder, file3))
	if err != nil {
		panic(err)
	}
	data4, err := os.ReadFile(path.Join(f.romFolder, file4))
	if err != nil {
		panic(err)
	}

	result := make([]uint8, len(data1)+len(data2)+len(data3)+len(data4))

	// Take alternating bytes from each file to assemble the complete buffer
	offset := 0
	for x := 0; x < len(data1); x += 2 {
		result[offset] = data1[x]
		result[offset+1] = data1[x+1]
		result[offset+2] = data2[x]
		result[offset+3] = data2[x+1]
		result[offset+4] = data3[x]
		result[offset+5] = data3[x+1]
		result[offset+6] = data4[x]
		result[offset+7] = data4[x+1]
		offset += 8
	}

	return memory.LoadData(address, result)
}

func (f *finalfight) loadFileSwapUint16(memory memory.Memory[uint32], file string, address uint32) error {
	data, err := os.ReadFile(path.Join(f.romFolder, file))
	if err != nil {
		panic(err)
	}

	// Swap each pair of bytes
	for x := 0; x < len(data); x += 2 {
		first := data[x]
		data[x] = data[x+1]
		data[x+1] = first
	}

	return memory.LoadData(address, data)
}
