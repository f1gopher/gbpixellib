package cartridge

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/f1gopher/gbpixellib/log"
)

type systemMemory interface {
	LoadData(offset uint32, data []uint8) error
}

type cartridge interface {
	init(romFolder string)
	isValid() error
	load(memory systemMemory) error
}

func LoadGame(log *log.Log, romFolder string, memory systemMemory) error {

	folderName := filepath.Base(romFolder)

	cart := findCart(folderName)

	cart.init(romFolder)

	if err := cart.isValid(); err != nil {
		return err
	}

	return cart.load(memory)
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
