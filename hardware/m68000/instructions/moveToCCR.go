package instructions

import "fmt"

type moveToCCR struct {
	length int
	src    addressModeInfo
}

func CreateMOVEToCCR(value uint16, code []uint8, codeOffset int) *moveToCCR {
	src, _, instSize := readEffectiveAddressMode(value, 0, wordSize, code, codeOffset, codeOffset)
	return &moveToCCR{
		length: 2 + instSize,
		src:    src,
	}
}

func (m *moveToCCR) Name() string {
	return fmt.Sprintf(
		"move    %s, CCR",
		m.src.string())
}

func (m *moveToCCR) Length() int {
	return m.length
}
