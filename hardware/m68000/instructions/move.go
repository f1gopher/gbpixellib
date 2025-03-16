package instructions

import "fmt"

type move struct {
	src  addressModeInfo
	dest addressModeInfo

	length      instSize
	sizeInBytes int
}

func CreateMOVE(value uint16, code []uint8, codeOffset int) *move {
	readOffset := codeOffset
	size := moveSize(value, 12)
	src, readOffset, _ := readSrcAddressMode(value, 0, size, code, readOffset, codeOffset)
	dest, readOffset, _ := readDestAddressMode(value, 6, size, code, readOffset, codeOffset)

	if dest.mode == addressRegisterDirect ||
		dest.mode == programCounterIndirectWithDisplacement ||
		dest.mode == preIndexedImediate ||
		dest.mode == programCounterIndirectWithIndex8bitDisplacement ||
		size == zeroSize {
		return nil
	}
	if size == byteSize && src.mode == addressRegisterDirect {
		return nil
	}

	sizeInBytes := 2 + (readOffset - codeOffset)

	return &move{
		src:         src,
		dest:        dest,
		length:      size,
		sizeInBytes: sizeInBytes,
	}
}

func (m *move) Name() string {
	return fmt.Sprintf(
		"move.%s  %s, %s",
		m.length.string(),
		m.src.string(),
		m.dest.string())
}

func (m *move) Length() int {
	return m.sizeInBytes
}
