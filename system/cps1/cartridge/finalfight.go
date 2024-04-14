package cartridge

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

func (f *finalfight) load() error {

}
