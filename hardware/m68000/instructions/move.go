package instructions

import "fmt"

type move struct {
	size     instSize
	destMode addressMode
	destReg  byte
	srcMode  addressMode
	srcReg   byte

	srcData  uint32
	destData uint32

	sizeInBytes int
}

func CreateMOVE(value uint16, code []uint8, offset int) *move {
	size := moveSize(value, 12)
	destMode, destReg := destination(value, 6, size)
	srcMode, srcReg := source(value, 0, size)
	sizeInBytes := 2

	var srcData, destData uint32
	destDataOffset := offset

	if srcMode == absoluteShort ||
		(srcMode == preIndexedImediate && size == wordSize) ||
		srcMode == addressRegisterIndirectWithDisplacement {
		srcData = readData(code, offset, wordSize)
		destDataOffset += wordSize.bytesUsed()
		sizeInBytes += 2
	}

	if srcMode == absoluteLong ||
		(srcMode == preIndexedImediate && size == longSize) {
		srcData = readData(code, offset, longSize)
		destDataOffset += int(longSize.bytesUsed())
		sizeInBytes += 4
	}

	if destMode == absoluteShort ||
		destMode == addressRegisterIndirectWithDisplacement {
		destData = readData(code, destDataOffset, wordSize)
		sizeInBytes += 2
	}

	if destMode == absoluteLong {
		destData = readData(code, destDataOffset, longSize)
		sizeInBytes += 4
	}

	return &move{
		size:        size,
		destMode:    destMode,
		destReg:     destReg,
		srcMode:     srcMode,
		srcReg:      srcReg,
		srcData:     srcData,
		destData:    destData,
		sizeInBytes: sizeInBytes,
	}
}

func (m *move) Name() string {
	return fmt.Sprintf(
		"move.%s  %s, %s",
		m.size.string(),
		addressString(m.srcMode, m.srcReg, m.srcData, 0),
		addressString(m.destMode, m.destReg, m.destData, 0))
}

func (m *move) Length() int {
	return m.sizeInBytes
}
