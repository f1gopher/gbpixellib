package cartridge

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/f1gopher/gbpixellib/hardware/memory"
	"github.com/f1gopher/gbpixellib/log"
)

type cartridge interface {
	init(romFolder string)
	isValid() error
	load(_68000 memory.Memory[uint32], gfx memory.Memory[uint32]) error
}

func LoadGame(
	log *log.Log,
	romFolder string,
	_68000Memory memory.Memory[uint32],
	gfxMemory memory.Memory[uint32]) error {

	folderName := filepath.Base(romFolder)

	cart := findCart(folderName)

	cart.init(romFolder)

	if err := cart.isValid(); err != nil {
		return err
	}

	return cart.load(_68000Memory, gfxMemory)
}

func findCart(name string) cartridge {
	name = strings.ToLower(name)

	switch name {
	case "ffight":
		return &finalfight{}
	default:
		panic(fmt.Sprintf("Unknown game '%s'", name))
	}
}
