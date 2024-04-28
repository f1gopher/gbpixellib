package cartridge

import (
	"path/filepath"
	"strings"
)

type systemMemory interface {
	Load(offset uint16, data []uint8) error
}

type cartridge interface {
	init(romFolder string)
	isValid() error
	load(memory systemMemory) error
}

func LoadGame(romFolder string, memory systemMemory) error {

	parts := filepath.SplitList(romFolder)
	folderName := parts[len(parts)-1]

	cart := findCart(folderName)

	cart.init(folderName)

	if err := cart.isValid(romFolder); err != nil {
		return err
	}

	return cart.load(memory)
}

func findCart(name string) cartridge {
	name = strings.ToLower(name)

	switch name {
	case "ffight":
		return &finalfight{}
	}
}
