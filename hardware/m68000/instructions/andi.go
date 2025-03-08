package instructions

import "fmt"

type andi struct {
	length int
	size   instSize
	src    addressModeInfo
	data   uint32
}

func CreateANDI(value uint16, code []uint8, codeOffset int) *andi {
	size := twoBitSize(value, 6)
	src, _, instSize := readEffectiveAddressMode(value, 0, size, code, codeOffset, codeOffset)

	return &andi{
		length: 2 + instSize + size.bytesUsed(),
		size:   size,
		src:    src,
		data:   readData(code, codeOffset, size),
	}
}

func (a *andi) Name() string {
	return fmt.Sprintf(
		"andi.%s  #$%d, %s",
		a.size.string(),
		a.data,
		a.src.string(),
	)
}

func (a *andi) Length() int {
	return a.length
}
