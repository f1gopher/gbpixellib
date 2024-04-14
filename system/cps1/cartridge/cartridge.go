package cartridge

import (
	"path/filepath"
	"strings"
)

type cartridge interface {
	init(romFolder string)
	isValid() error
	load() error
}

func LoadGame(romFolder string) error {

	parts := filepath.SplitList(romFolder)
	folderName := parts[len(parts)-1]

	cart := findCart(folderName)

	cart.init(folderName)

	if err := cart.isValid(romFolder); err != nil {

	}

	return cart.load(romFolder)
}

func findCart(name string) cartridge {
	name = strings.ToLower(name)

	switch name {
	case "ffight":
		return &finalfight{}
	}
}
